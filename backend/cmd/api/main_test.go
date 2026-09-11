package main

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadConfigDefaultsToSecureCookieAndThirtyMinutes(t *testing.T) {
	t.Setenv("RPMP_DATABASE_URL", "postgres://example")
	t.Setenv("RPMP_JWT_SECRET", "01234567890123456789012345678901")
	t.Setenv("RPMP_ALLOWED_ORIGIN", "https://rpmp.example")
	t.Setenv("RPMP_ACCESS_TTL", "")
	t.Setenv("RPMP_LOCAL_HTTP", "")
	t.Setenv("RPMP_SOURCE_CSV_PATH", "/tmp/source.csv")

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
	t.Setenv("RPMP_SOURCE_CSV_PATH", "/tmp/source.csv")

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
	t.Setenv("RPMP_SOURCE_CSV_PATH", "/tmp/source.csv")
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
	t.Setenv("RPMP_SOURCE_CSV_PATH", "")
	contents := "RPMP_DATABASE_URL=postgres://example\n" +
		"RPMP_JWT_SECRET=01234567890123456789012345678901\n" +
		"RPMP_ALLOWED_ORIGIN=http://localhost:3000\n" +
		"RPMP_LOCAL_HTTP=true\n" +
		"RPMP_SOURCE_CSV_PATH=/tmp/source.csv\n"
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
	t.Setenv("RPMP_SOURCE_CSV_PATH", "/tmp/source.csv")
	if _, err := loadConfig(); err == nil {
		t.Fatal("insecure non-local origin accepted")
	}
}

func TestLogStartupOmitsSecrets(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	secret := "01234567890123456789012345678901"
	logStartup(logger, config{
		addr:          ":8080",
		databaseURL:   "postgres://user:" + secret + "@localhost:5432/rpmp",
		jwtSecret:     secret,
		allowedOrigin: "http://localhost:3000",
		accessTTL:     30 * time.Minute,
		localHTTP:     true,
	})
	line := buf.String()
	if !strings.Contains(line, `msg="API listening"`) || !strings.Contains(line, "addr=:8080") {
		t.Fatalf("missing listen fields: %s", line)
	}
	if !strings.Contains(line, "origin=http://localhost:3000") || !strings.Contains(line, "local_http=true") {
		t.Fatalf("missing origin fields: %s", line)
	}
	if !strings.Contains(line, "access_ttl=") {
		t.Fatalf("missing ttl: %s", line)
	}
	if strings.Contains(line, secret) || strings.Contains(line, "postgres://") || strings.Contains(line, "RPMP_JWT_SECRET") {
		t.Fatalf("startup log leaked a secret: %s", line)
	}
}
