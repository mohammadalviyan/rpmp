package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mohammadalviyan/rpmp/backend/internal/adapter/stub"
	"github.com/mohammadalviyan/rpmp/backend/internal/auth"
	"github.com/mohammadalviyan/rpmp/backend/internal/handler"
	"github.com/mohammadalviyan/rpmp/backend/internal/repo"
	"github.com/mohammadalviyan/rpmp/backend/internal/usecase"
)

type config struct {
	addr          string
	databaseURL   string
	jwtIssuer     string
	jwtSecret     string
	allowedOrigin string
	accessTTL     time.Duration
	localHTTP     bool
}

func loadConfig() (config, error) {
	if err := loadDotEnv(".env"); err != nil {
		return config{}, err
	}
	cfg := config{
		addr:          envOrDefault("RPMP_ADDR", ":8080"),
		databaseURL:   os.Getenv("RPMP_DATABASE_URL"),
		jwtIssuer:     envOrDefault("RPMP_JWT_ISSUER", "rpmp"),
		jwtSecret:     os.Getenv("RPMP_JWT_SECRET"),
		allowedOrigin: os.Getenv("RPMP_ALLOWED_ORIGIN"),
		accessTTL:     30 * time.Minute,
	}
	if raw := os.Getenv("RPMP_ACCESS_TTL"); raw != "" {
		ttl, err := time.ParseDuration(raw)
		if err != nil || ttl <= 0 {
			return config{}, errors.New("RPMP_ACCESS_TTL must be a positive Go duration")
		}
		cfg.accessTTL = ttl
	}
	if raw := os.Getenv("RPMP_LOCAL_HTTP"); raw != "" {
		localHTTP, err := strconv.ParseBool(raw)
		if err != nil {
			return config{}, errors.New("RPMP_LOCAL_HTTP must be true or false")
		}
		cfg.localHTTP = localHTTP
	}
	if cfg.databaseURL == "" {
		return config{}, errors.New("RPMP_DATABASE_URL is required")
	}
	if cfg.jwtSecret == "" {
		return config{}, errors.New("RPMP_JWT_SECRET is required")
	}
	if len(cfg.jwtSecret) < 32 {
		return config{}, errors.New("RPMP_JWT_SECRET must contain at least 32 bytes")
	}
	if cfg.allowedOrigin == "" {
		return config{}, errors.New("RPMP_ALLOWED_ORIGIN is required")
	}
	origin, err := url.Parse(cfg.allowedOrigin)
	if err != nil || origin.Host == "" || origin.User != nil || origin.RawQuery != "" || origin.Fragment != "" || origin.Path != "" {
		return config{}, errors.New("RPMP_ALLOWED_ORIGIN must be an absolute HTTP origin")
	}
	if cfg.localHTTP {
		ip := net.ParseIP(origin.Hostname())
		if origin.Scheme != "http" || (origin.Hostname() != "localhost" && (ip == nil || !ip.IsLoopback())) {
			return config{}, errors.New("RPMP_LOCAL_HTTP requires a localhost or loopback HTTP origin")
		}
	} else if origin.Scheme != "https" {
		return config{}, errors.New("RPMP_ALLOWED_ORIGIN must use HTTPS outside local development")
	}
	return cfg, nil
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
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
		if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
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
		return fmt.Errorf("configure database pool: %w", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	logStartup(slog.Default(), cfg)

	passwords, err := auth.NewPasswords(0)
	if err != nil {
		return err
	}
	dummyHash, err := passwords.Hash("rpmp-invalid-" + time.Now().UTC().String())
	if err != nil {
		return fmt.Errorf("prepare credential verifier: %w", err)
	}
	tokens, err := auth.NewTokens(cfg.jwtIssuer, cfg.jwtSecret, cfg.accessTTL)
	if err != nil {
		return err
	}

	postgres := repo.NewPostgres(pool)
	login := usecase.NewLogin(postgres, postgres, passwords, tokens, dummyHash)
	current := usecase.NewCurrentUser(postgres)
	logout := usecase.NewLogout(postgres)
	source, err := stub.New()
	if err != nil {
		return fmt.Errorf("load source fixtures: %w", err)
	}
	dashboard := usecase.NewDashboardSummary(source)
	authHandler := handler.NewAuth(login, current, logout, handler.CookieConfig{
		TTL: cfg.accessTTL, Secure: !cfg.localHTTP,
	})
	dashboardHandler := handler.NewDashboard(dashboard)

	server := &http.Server{
		Addr:              cfg.addr,
		Handler:           handler.NewRouter(authHandler, dashboardHandler, tokens, cfg.allowedOrigin),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	errs := make(chan error, 1)
	go func() {
		errs <- server.ListenAndServe()
	}()
	select {
	case err := <-errs:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("API listen failed", "error", err)
			return err
		}
	case <-ctx.Done():
		slog.Info("API shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
	return nil
}

func logStartup(logger *slog.Logger, cfg config) {
	logger.Info("API listening",
		"addr", cfg.addr,
		"origin", cfg.allowedOrigin,
		"local_http", cfg.localHTTP,
		"access_ttl", cfg.accessTTL,
	)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx); err != nil {
		slog.Error("API stopped", "error", err)
		os.Exit(1)
	}
}
