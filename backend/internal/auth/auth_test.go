package auth_test

import (
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