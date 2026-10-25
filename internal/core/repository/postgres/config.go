package postgres

import "time"

type postgresConfig struct {
	connStr           string
	maxConns          int32
	minConns          int32
	maxConnLifetime   time.Duration
	maxConnIdleTime   time.Duration
	healthCheckPeriod time.Duration
}

func NewPostgresConfig(connStr string, maxConns, minConns int32, connMaxLifetime,
	connMaxIdleTime, healthCheckPeriod time.Duration) *postgresConfig {
	return &postgresConfig{
		connStr:           connStr,
		maxConns:          maxConns,
		minConns:          minConns,
		maxConnLifetime:   connMaxLifetime,
		maxConnIdleTime:   connMaxIdleTime,
		healthCheckPeriod: healthCheckPeriod,
	}
}
