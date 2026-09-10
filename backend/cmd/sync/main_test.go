package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigUsesCSVModeWithoutSourceDatabaseURL(t *testing.T) {
	t.Setenv("RPMP_DATABASE_URL", "postgres://rpmp")
	t.Setenv("RPMP_SOURCE_CSV_PATH", "../source.csv")
	t.Setenv("RPMP_SOURCE_DATABASE_URL", "")

	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.databaseURL != "postgres://rpmp" || cfg.csvPath != "../source.csv" {
		t.Fatalf("config = %#v", cfg)
	}
}

func TestLoadConfigReadsDotEnv(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("RPMP_DATABASE_URL", "")
	t.Setenv("RPMP_SOURCE_CSV_PATH", "")
	if err := os.WriteFile(
		filepath.Join(dir, ".env"),
		[]byte("RPMP_DATABASE_URL=postgres://rpmp\nRPMP_SOURCE_CSV_PATH=../source.csv\n"),
		0o600,
	); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.databaseURL != "postgres://rpmp" || cfg.csvPath != "../source.csv" {
		t.Fatalf("dotenv config = %#v", cfg)
	}
}

func TestLoadConfigRequiresCSVPath(t *testing.T) {
	t.Setenv("RPMP_DATABASE_URL", "postgres://rpmp")
	t.Setenv("RPMP_SOURCE_CSV_PATH", "")
	if _, err := loadConfig(); err == nil {
		t.Fatal("missing CSV path accepted")
	}
}
