package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type storage struct {
	conn *pgxpool.Pool
}

func newStorage(conn *pgxpool.Pool) *storage {
	return &storage{conn: conn}
}

type test struct {
	ID    string
	Title string
}

func (s *storage) first() (test, error) {
	ctx := context.Background()

	first := test{}
	// Uncomment for retries:
	//err := AcquireRetryOnceFunc(ctx, s.conn, func(conn *pgxpool.Conn) error {
	err := s.conn.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		row := conn.QueryRow(ctx, "SELECT id, title FROM test LIMIT 1")
		if errScan := row.Scan(&first.ID, &first.Title); errScan != nil {
			return errScan
		}
		return nil
	})
	if err != nil {
		return first, err
	}
	return first, nil
}

func (s *storage) list() ([]test, error) {
	ctx := context.Background()

	tests := make([]test, 0)

	// Uncomment for retries:
	//err := AcquireRetryOnceFunc(ctx, s.conn, func(conn *pgxpool.Conn) error {
	err := s.conn.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		rows, errQuery := conn.Query(ctx, "SELECT id, title FROM test LIMIT 100")
		defer rows.Close()

		if errQuery != nil {
			return errQuery
		}

		for rows.Next() {
			var row test
			if errScan := rows.Scan(&row.ID, &row.Title); errScan != nil {
				return errScan
			}
			tests = append(tests, row)
		}

		return rows.Err()
	})
	if err != nil {
		return tests, err
	}
	return tests, nil
}

func AcquireRetryOnceFunc(ctx context.Context, db *pgxpool.Pool, f func(*pgxpool.Conn) error) error {
	err := db.AcquireFunc(ctx, f)
	if err != nil {
		err = db.AcquireFunc(ctx, f)
	}
	return err
}
