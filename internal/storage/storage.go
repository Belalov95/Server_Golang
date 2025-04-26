package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func GetConnect(connStr string) (*pgx.Conn, error) { //создается для соединения с бд
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		wrappedErr := errors.Wrap(err, "failed to connect to the database")
		log.Fatalf("Unable to connect to database: %v\n", wrappedErr)
		return nil, wrappedErr
	}
	return conn, nil
}
