package main

import (
	"os"
	"path/filepath"
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

func TestLoadConfigReadsDotEnvWhenProcessEnvIsUnset(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("RPMP_DATABASE_URL", "")
	t.Setenv("RPMP_JWT_SECRET", "")
	t.Setenv("RPMP_ALLOWED_ORIGIN", "")
	t.Setenv("RPMP_LOCAL_HTTP", "")
	contents := "RPMP_DATABASE_URL=postgres://example\n" +
		"RPMP_JWT_SECRET=01234567890123456789012345678901\n" +
		"RPMP_ALLOWED_ORIGIN=http://localhost:3000\n" +
		"RPMP_LOCAL_HTTP=true\n"
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.databaseURL != "postgres://example" || !cfg.localHTTP {
		t.Fatalf("dotenv was not applied: %#v", cfg)
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
