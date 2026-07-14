package jwt

import (
	"testing"
	"time"
)

func TestGenerateAndVerifyToken(t *testing.T) {
	token, err := GenerateToken("user-1", "bayu@example.com", "secret", time.Hour, "portfolio")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	claims, err := VerifyToken(token, "secret")
	if err != nil {
		t.Fatalf("VerifyToken returned error: %v", err)
	}

	if claims.UserID != "user-1" || claims.Email != "bayu@example.com" || claims.Issuer != "portfolio" {
		t.Fatalf("claims were not preserved: %+v", claims)
	}
}

func TestVerifyTokenRejectsInvalidTokens(t *testing.T) {
	token, err := GenerateToken("user-1", "bayu@example.com", "secret", time.Hour, "portfolio")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	if _, err := VerifyToken(token, "different-secret"); err == nil {
		t.Fatal("expected invalid signature error")
	}

	if _, err := VerifyToken("not-a-token", "secret"); err == nil {
		t.Fatal("expected malformed token error")
	}
}
