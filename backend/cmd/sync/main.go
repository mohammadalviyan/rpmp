package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mohammadalviyan/rpmp/backend/internal/adapter/csvsource"
	"github.com/mohammadalviyan/rpmp/backend/internal/repo"
	"github.com/mohammadalviyan/rpmp/backend/internal/usecase"
)

type config struct {
	databaseURL string
	csvPath     string
}

func loadConfig() (config, error) {
	if err := loadDotEnv(".env"); err != nil {
		return config{}, err
	}
	cfg := config{
		databaseURL: os.Getenv("RPMP_DATABASE_URL"),
		csvPath:     os.Getenv("RPMP_SOURCE_CSV_PATH"),
	}
	if cfg.databaseURL == "" {
		return config{}, errors.New("RPMP_DATABASE_URL is required")
	}
	if cfg.csvPath == "" {
		return config{}, errors.New("RPMP_SOURCE_CSV_PATH is required for CSV mode")
	}
	return cfg, nil
}

func loadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("%s:%d: expected KEY=VALUE", path, lineNumber)
		}
		key = strings.TrimSpace(key)
		if key == "" {
			return fmt.Errorf("%s:%d: empty key", path, lineNumber)
		}
		if os.Getenv(key) != "" {
			continue
		}
		if err := os.Setenv(key, unquoteEnv(strings.TrimSpace(value))); err != nil {
			return fmt.Errorf("set %s: %w", key, err)
		}
	}
	return scanner.Err()
}

func unquoteEnv(value string) string {
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			return value[1 : len(value)-1]
		}
	}
	return value
}

func run(ctx context.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	pool, err := pgxpool.New(ctx, cfg.databaseURL)
	if err != nil {
		return fmt.Errorf("configure RPMP database pool: %w", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("connect to RPMP database: %w", err)
	}

	runner := usecase.NewSyncRunner(csvsource.New(cfg.csvPath), repo.NewSyncPostgres(pool))
	result, err := runner.Run(ctx)
	if err != nil {
		return err
	}
	slog.Info("sync completed",
		"rows_read", result.Run.RowsRead,
		"rows_written", result.Run.RowsWritten,
		"idempotent", result.Idempotent,
	)
	return nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx); err != nil {
		slog.Error("sync failed", "error", err)
		os.Exit(1)
	}
}
