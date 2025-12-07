package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/GeorgeTyupin/labguard/internal/server/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MustDBPoolInit(logger *slog.Logger, pgConf config.PostgresConfig) *pgxpool.Pool {
	const op = "server.repository.postgres.MustInit"
	logger = logger.With(slog.String("op", op))

	pool, err := newPool(pgConf)
	if err != nil {
		logger.Error("Не удалось подключиться к бд", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return pool
}

func newPool(pgConf config.PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		pgConf.Postgres.User,
		pgConf.Postgres.Password,
		pgConf.Postgres.Host,
		pgConf.Postgres.Port,
		pgConf.Postgres.Database,
	)

	poolConf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить DSN: %w", err)
	}

	poolConf.MaxConns = pgConf.Postgres.PoolSize
	poolConf.MaxConnLifetime = pgConf.Postgres.Connection.MaxLifeTime
	poolConf.MaxConnIdleTime = pgConf.Postgres.Connection.MaxIdleTime
	poolConf.HealthCheckPeriod = pgConf.Postgres.Connection.HealthCheckPeriod
	poolConf.ConnConfig.ConnectTimeout = pgConf.Postgres.Connection.Timeout

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConf)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать пулл соединений к базе данных: %w", err)
	}

	return pool, nil
}
