package storage

import (
	"context"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func GetConnect(ctx context.Context, connStr string) (*pgxpool.Pool, error) { // создается для соединения с бд
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	// Добавляем трассировщик от otelpgx
	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create pgxpool")
	}

	// проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "Unable to ping database")
	}

	return pool, nil
}
