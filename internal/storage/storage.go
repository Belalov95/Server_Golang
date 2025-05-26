package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func GetConnect(connStr string) (*pgx.Conn, error) { //создается для соединения с бд
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}
	return conn, nil
}
