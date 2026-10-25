package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePool(ctx context.Context, cfg *postgresConfig, logger *slog.Logger) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.connStr)
	if err != nil {
		logger.Error("unable to parse DSN", slog.Any("err", err))
		return nil, fmt.Errorf("unable to parse DSN: %w", err)
	}

	poolConfig.MaxConns = cfg.maxConns
	poolConfig.MinConns = cfg.minConns
	poolConfig.MaxConnLifetime = cfg.maxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.maxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.healthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Error("failed to create connection pool", slog.Any("err", err))
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Error("unable to ping database", slog.Any("err", err))
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	logger.Info("database connection pool created and verified",
		slog.Int("max_conns", int(cfg.maxConns)),
		slog.Int("min_conns", int(cfg.minConns)),
	)

	return pool, nil
}

func GetPostgresValues() string {
	url := os.Getenv("POSTGRES_URL")
	if url != "" {
		return url
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		password = "test123"
	}

	dbname := os.Getenv("POSTGRES_DB")
	if dbname == "" {
		dbname = "habit-graph_app"
	}

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5433"
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
}
