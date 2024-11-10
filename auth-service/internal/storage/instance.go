package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vlks-dev/sibavto/shared/utils/config"
	"log/slog"
	"time"
)

var (
	MaxConnIdleTimeErr = "failed to parse max connection idle time"
)

func PostgresPool(ctx context.Context, config *config.Config, log *slog.Logger) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	postgres := config.Database.PostgresSQL
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		postgres.Username,
		postgres.Password,
		postgres.Host,
		postgres.Port,
		postgres.Database,
	)

	configdb, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Error("failed to parse postgres connection url", "err", err)
		return nil, err
	}

	configdb.MaxConns = postgres.MaxConnections
	configdb.MinConns = postgres.MinConnections
	configdb.MaxConnIdleTime, err = time.ParseDuration(postgres.IdleTime)
	if err != nil {
		log.Error("failed to parse postgres connection idle time", "err", err, "package")
		return nil, errors.New(MaxConnIdleTimeErr)
	}
	configdb.HealthCheckPeriod, err = time.ParseDuration(postgres.HealthCheckPeriod)
	if err != nil {
		log.Error("failed to parse postgres health check period", "err", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, configdb)
	if err != nil {
		log.Error("failed to create postgres connection pool", "err", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		log.Error("failed to ping postgres", "url", url)
		return nil, err
	}

	log.Debug(
		"postgres pool connection established",
		"max connections",
		configdb.MaxConns,
		"max connections idle",
		configdb.MaxConnIdleTime.Seconds(),
		"health check period",
		configdb.HealthCheckPeriod.Seconds(),
	)

	return pool, nil
}
