package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	secret := "mysecret"
	userID := uuid.New()
	expiresIn := time.Hour

	// Test token creation
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Test token validation
	returnedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}
	if returnedUserID != userID {
		t.Fatalf("Expected user ID %v, got %v", userID, returnedUserID)
	}

	// Test expired token
	expiredToken, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("Failed to create expired JWT: %v", err)
	}
	_, err = ValidateJWT(expiredToken, secret)
	if err == nil {
		t.Fatal("Expected error for expired token, got none")
	}

	// Test invalid token
	_, err = ValidateJWT("invalidtoken", secret)
	if err == nil {
		t.Fatal("Expected error for invalid token, got none")
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := make(map[string][]string)
	headers["Authorization"] = []string{"Bearer mytoken"}
	token, err:= GetBearerToken(headers)
	if err != nil {
		t.Fatalf("Failed to get bearer token: %v", err)
	}
	if token != "mytoken" {
		t.Fatalf("Expected token 'mytoken', got '%s'", token)
	}

	// Test missing header
	headers = make(map[string][]string)
	_, err = GetBearerToken(headers)
	if err == nil {
		t.Fatal("Expected error for missing Authorization header, got none")
	}
}

func TestHashRefreshToken(t *testing.T) {
	token := "test-refresh-token"

	got := HashRefreshToken(token)

	expectedHash := sha256.Sum256([]byte(token))
	expected := hex.EncodeToString(expectedHash[:])

	if got != expected {
		t.Fatalf("Expected hash %q, got %q", expected, got)
	}

	if got == token {
		t.Fatal("Refresh token was stored as plaintext")
	}

	if len(got) != 64 {
		t.Fatalf("Expected SHA-256 hex hash to be 64 characters, got %d", len(got))
	}
}

func TestHashRefreshTokenIsDeterministic(t *testing.T) {
	token := "same-refresh-token"

	first := HashRefreshToken(token)
	second := HashRefreshToken(token)

	if first != second {
		t.Fatalf(
			"Expected same token to produce same hash, got %q and %q",
			first,
			second,
		)
	}
}

func TestHashRefreshTokenDifferentTokens(t *testing.T) {
	first := HashRefreshToken("refresh-token-1")
	second := HashRefreshToken("refresh-token-2")

	if first == second {
		t.Fatal("Different refresh tokens produced the same hash")
	}
}

func TestJWTWithDifferentUsers(t *testing.T) {
	secret := "mysecret"
	firstUserID := uuid.New()
	secondUserID := uuid.New()

	firstToken, err := MakeJWT(firstUserID, secret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create first JWT: %v", err)
	}

	secondToken, err := MakeJWT(secondUserID, secret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create second JWT: %v", err)
	}

	returnedFirstID, err := ValidateJWT(firstToken, secret)
	if err != nil {
		t.Fatalf("Failed to validate first JWT: %v", err)
	}

	returnedSecondID, err := ValidateJWT(secondToken, secret)
	if err != nil {
		t.Fatalf("Failed to validate second JWT: %v", err)
	}

	if returnedFirstID != firstUserID {
		t.Fatalf(
			"Expected first user ID %v, got %v",
			firstUserID,
			returnedFirstID,
		)
	}

	if returnedSecondID != secondUserID {
		t.Fatalf(
			"Expected second user ID %v, got %v",
			secondUserID,
			returnedSecondID,
		)
	}
}

func TestJWTWrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, "correct-secret", time.Hour)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Fatal("Expected validation to fail with wrong secret")
	}
}

func TestGetBearerTokenInvalidFormat(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{
			name:   "missing bearer prefix",
			header: "mytoken",
		},
		{
			name:   "empty bearer token",
			header: "Bearer ",
		},
		{
			name:   "wrong authentication scheme",
			header: "Basic mytoken",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := map[string][]string{
				"Authorization": {tt.header},
			}

			_, err := GetBearerToken(headers)
			if err == nil {
				t.Fatalf(
					"Expected error for Authorization header %q",
					tt.header,
				)
			}
		})
	}
}

func TestHashRefreshTokenEmptyToken(t *testing.T) {
	got := HashRefreshToken("")

	expectedHash := sha256.Sum256([]byte(""))
	expected := hex.EncodeToString(expectedHash[:])

	if got != expected {
		t.Fatalf(
			"Expected hash %q, got %q",
			expected,
			got,
		)
	}
}

func TestHashRefreshTokenChangesWhenInputChanges(t *testing.T) {
	first := HashRefreshToken("refresh-token")
	second := HashRefreshToken("refresh-token-changed")

	if first == second {
		t.Fatal("Changing the refresh token should change its hash")
	}
}

func TestValidateRegistration(t *testing.T) {
	tests := []struct {
		name        string
		request     register
		wantError   bool
		wantEmail   string
		wantUsername string
	}{
		{
			name: "valid registration",
			request: register{
				Email:    "user@example.com",
				Username: "testuser",
				Password: "password123",
			},
			wantError:    false,
			wantEmail:    "user@example.com",
			wantUsername: "testuser",
		},
		{
			name: "trims email",
			request: register{
				Email:    "  user@example.com  ",
				Username: "testuser",
				Password: "password123",
			},
			wantError:    false,
			wantEmail:    "user@example.com",
			wantUsername: "testuser",
		},
		{
			name: "missing email",
			request: register{
				Username: "testuser",
				Password: "password123",
			},
			wantError: true,
		},
		{
			name: "invalid email",
			request: register{
				Email:    "not-an-email",
				Username: "testuser",
				Password: "password123",
			},
			wantError: true,
		},
		{
			name: "password too short",
			request: register{
				Email:    "user@example.com",
				Username: "testuser",
				Password: "1234567",
			},
			wantError: true,
		},
		{
			name: "missing username",
			request: register{
				Email:    "user@example.com",
				Password: "password123",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.request

			err := validateRegistration(&req)

			if tt.wantError && err == nil {
				t.Fatal("Expected validation error, got nil")
			}

			if !tt.wantError && err != nil {
				t.Fatalf(
					"Expected no validation error, got: %v",
					err,
				)
			}

			if !tt.wantError {
				if req.Email != tt.wantEmail {
					t.Fatalf(
						"Expected email %q, got %q",
						tt.wantEmail,
						req.Email,
					)
				}

				if req.Username != tt.wantUsername {
					t.Fatalf(
						"Expected username %q, got %q",
						tt.wantUsername,
						req.Username,
					)
				}
			}
		})
	}
}