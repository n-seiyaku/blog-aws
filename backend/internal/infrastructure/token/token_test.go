package token

import (
	"testing"
	"time"
)

func TestGenerateRefreshTokenAndHash(t *testing.T) {
	firstToken, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken returned error: %v", err)
	}
	secondToken, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken returned error: %v", err)
	}

	if firstToken == "" || secondToken == "" {
		t.Fatal("expected non-empty refresh tokens")
	}
	if firstToken == secondToken {
		t.Fatal("expected refresh tokens to be unique")
	}
	// ハッシュ処理の決定性を検証
	hash1 := HashRefreshToken(firstToken)
	hash2 := HashRefreshToken(firstToken)
	if hash1 != hash2 {
		t.Fatal("expected refresh token hash to be deterministic")
	}
	if hash1 == firstToken {
		t.Fatal("expected refresh token hash to differ from raw token")
	}
}

func TestJWTGeneratorGenerateAndVerifyToken(t *testing.T) {
	generator := NewJWTGenerator("test-secret", time.Hour)

	tokenString, err := generator.GenerateToken("user-1", "user@example.com", "ADMIN")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	claims, err := generator.VerifyToken(tokenString)
	if err != nil {
		t.Fatalf("VerifyToken returned error: %v", err)
	}

	if claims.Subject != "user-1" {
		t.Fatalf("expected subject user-1, got %q", claims.Subject)
	}
	if claims.Email != "user@example.com" {
		t.Fatalf("expected email user@example.com, got %q", claims.Email)
	}
	if claims.Role != "ADMIN" {
		t.Fatalf("expected role ADMIN, got %q", claims.Role)
	}
}

func TestJWTGeneratorRejectsInvalidToken(t *testing.T) {
	generator := NewJWTGenerator("test-secret", time.Hour)

	if _, err := generator.VerifyToken("not-a-token"); err == nil {
		t.Fatal("expected invalid token error")
	}
}
