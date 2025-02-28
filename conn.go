package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newConn(connString string) *pgxpool.Pool {
	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connection to database: %v\n", err))
	}

	return conn
}
