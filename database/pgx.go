package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func PGXConnect(dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func PGXExec(db *pgx.Conn, ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	if db == nil {
		return pgconn.CommandTag{}, errors.New("database connection is not initialized")
	}
	return db.Exec(ctx, query, args...)
}

func PGXExecNoResult(db *pgx.Conn, ctx context.Context, query string, args ...interface{}) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}
	_, err := PGXExec(db, ctx, query, args...)
	return err
}

func PGXQuery(db *pgx.Conn, ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	return db.Query(ctx, query, args...)
}

func PGXQueryRow(db *pgx.Conn, ctx context.Context, query string, args ...interface{}) pgx.Row {
	if db == nil {
		return nil
	}
	return db.QueryRow(ctx, query, args...)
}
