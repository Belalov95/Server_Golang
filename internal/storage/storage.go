package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func GetConnect(connStr string) (*pgx.Conn, error) { //создается для соединения с бд
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	return conn, nil
}
