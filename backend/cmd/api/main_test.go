package main

import (
	"testing"
	"time"
)

func TestLoadConfigDefaultsToSecureCookieAndThirtyMinutes(t *testing.T) {
	t.Setenv("RPMP_DATABASE_URL", "postgres://example")
	t.Setenv("RPMP_JWT_SECRET", "01234567890123456789012345678901")
	t.Setenv("RPMP_ALLOWED_ORIGIN", "https://rpmp.example")
	t.Setenv("RPMP_ACCESS_TTL", "")
	t.Setenv("RPMP_LOCAL_HTTP", "")

	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.accessTTL != 30*time.Minute || cfg.localHTTP {
		t.Fatalf("unexpected defaults: %#v", cfg)
	}
}

func TestLoadConfigAllowsExplicitLocalHTTP(t *testing.T) {
	t.Setenv("RPMP_DATABASE_URL", "postgres://example")
	t.Setenv("RPMP_JWT_SECRET", "01234567890123456789012345678901")
	t.Setenv("RPMP_ALLOWED_ORIGIN", "http://localhost:3000")
	t.Setenv("RPMP_LOCAL_HTTP", "true")

	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.localHTTP {
		t.Fatal("explicit local HTTP was not enabled")
	}
}

func TestLoadConfigRequiresServerSecrets(t *testing.T) {
	t.Setenv("RPMP_DATABASE_URL", "postgres://example")
	t.Setenv("RPMP_JWT_SECRET", "")
	t.Setenv("RPMP_ALLOWED_ORIGIN", "https://rpmp.example")
	if _, err := loadConfig(); err == nil {
		t.Fatal("missing JWT secret accepted")
	}
}

func TestLoadConfigRejectsInsecureNonLocalOrigin(t *testing.T) {
	t.Setenv("RPMP_DATABASE_URL", "postgres://example")
	t.Setenv("RPMP_JWT_SECRET", "01234567890123456789012345678901")
	t.Setenv("RPMP_ALLOWED_ORIGIN", "http://rpmp.example")
	t.Setenv("RPMP_LOCAL_HTTP", "true")
	if _, err := loadConfig(); err == nil {
		t.Fatal("insecure non-local origin accepted")
	}
}
