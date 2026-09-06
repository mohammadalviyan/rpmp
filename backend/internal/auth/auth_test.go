package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

const testKey = "01234567890123456789012345678901"

func TestPasswordsHashAndVerify(t *testing.T) {
	passwords, err := NewPasswords(4)
	if err != nil {
		t.Fatal(err)
	}
	hash, err := passwords.Hash("correct horse battery staple")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(hash, "correct horse") {
		t.Fatal("hash contains plaintext password")
	}
	if err := passwords.Verify(hash, "correct horse battery staple"); err != nil {
		t.Fatalf("correct password rejected: %v", err)
	}
	if err := passwords.Verify(hash, "wrong"); err == nil {
		t.Fatal("wrong password accepted")
	}
}

func TestTokensIssueAndVerify(t *testing.T) {
	tokens, err := NewTokens("rpmp-test", testKey, 30*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 9, 6, 7, 0, 0, 0, time.UTC)
	tokens.now = func() time.Time { return now }
	tokens.newJTI = func() string { return "jti-1" }

	raw, err := tokens.Issue("user-1", domain.RoleViewer)
	if err != nil {
		t.Fatal(err)
	}
	principal, err := tokens.Verify(raw)
	if err != nil {
		t.Fatal(err)
	}
	if principal.UserID != "user-1" || principal.Role != domain.RoleViewer {
		t.Fatalf("unexpected principal: %#v", principal)
	}

	parsed, _, err := jwt.NewParser().ParseUnverified(raw, &Claims{})
	if err != nil {
		t.Fatal(err)
	}
	claims := parsed.Claims.(*Claims)
	if claims.Issuer != "rpmp-test" || claims.Subject != "user-1" || claims.ID != "jti-1" {
		t.Fatalf("required claims missing: %#v", claims)
	}
}

func TestTokensRejectInvalidTokens(t *testing.T) {
	now := time.Date(2026, 9, 6, 7, 0, 0, 0, time.UTC)
	tokens, _ := NewTokens("rpmp-test", testKey, time.Minute)
	tokens.now = func() time.Time { return now }
	tokens.newJTI = func() string { return "jti-1" }

	valid, _ := tokens.Issue("user-1", domain.RoleAdmin)
	tokens.now = func() time.Time { return now.Add(2 * time.Minute) }
	if _, err := tokens.Verify(valid); domain.ErrorKindOf(err) != domain.KindUnauthenticated {
		t.Fatalf("expired token error = %v", err)
	}

	other, _ := NewTokens("rpmp-test", "abcdefghijklmnopqrstuvwxyz012345", time.Minute)
	other.now = func() time.Time { return now }
	other.newJTI = func() string { return "jti-2" }
	tampered, _ := other.Issue("user-1", domain.RoleAdmin)
	if _, err := tokens.Verify(tampered); domain.ErrorKindOf(err) != domain.KindUnauthenticated {
		t.Fatalf("wrong signature error = %v", err)
	}

	missingJTI := Claims{
		Role: domain.RoleViewer,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "rpmp-test", Subject: "user-1",
			IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
		},
	}
	raw, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, missingJTI).SignedString([]byte(testKey))
	tokens.now = func() time.Time { return now }
	if _, err := tokens.Verify(raw); domain.ErrorKindOf(err) != domain.KindUnauthenticated {
		t.Fatalf("missing claim error = %v", err)
	}
}

func TestTokensRejectWeakConfiguration(t *testing.T) {
	if _, err := NewTokens("rpmp", "short", time.Minute); err == nil {
		t.Fatal("weak key accepted")
	}
}
