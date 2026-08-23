package auth_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/EyuAtske/AfriMart/backend/internal/auth"
)

func TestJWT(t *testing.T) {
	secret := "mysecret"
	userID := uuid.New()
	expiresIn := time.Hour

	// Test token creation
	token, err := auth.MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Test token validation
	returnedUserID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}
	if returnedUserID != userID {
		t.Fatalf("Expected user ID %v, got %v", userID, returnedUserID)
	}

	// Test expired token
	expiredToken, err := auth.MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("Failed to create expired JWT: %v", err)
	}
	_, err = auth.ValidateJWT(expiredToken, secret)
	if err == nil {
		t.Fatal("Expected error for expired token, got none")
	}

	// Test invalid token
	_, err = auth.ValidateJWT("invalidtoken", secret)
	if err == nil {
		t.Fatal("Expected error for invalid token, got none")
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := make(map[string][]string)
	headers["Authorization"] = []string{"Bearer mytoken"}
	token, err:= auth.GetBearerToken(headers)
	if err != nil {
		t.Fatalf("Failed to get bearer token: %v", err)
	}
	if token != "mytoken" {
		t.Fatalf("Expected token 'mytoken', got '%s'", token)
	}

	// Test missing header
	headers = make(map[string][]string)
	_, err = auth.GetBearerToken(headers)
	if err == nil {
		t.Fatal("Expected error for missing Authorization header, got none")
	}
}

func TestHashRefreshToken(t *testing.T) {
	token := "test-refresh-token"

	got := auth.HashRefreshToken(token)

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

	first := auth.HashRefreshToken(token)
	second := auth.HashRefreshToken(token)

	if first != second {
		t.Fatalf(
			"Expected same token to produce same hash, got %q and %q",
			first,
			second,
		)
	}
}

func TestHashRefreshTokenDifferentTokens(t *testing.T) {
	first := auth.HashRefreshToken("refresh-token-1")
	second := auth.HashRefreshToken("refresh-token-2")

	if first == second {
		t.Fatal("Different refresh tokens produced the same hash")
	}
}