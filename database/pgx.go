package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	// "github.com/jackc/pgx/v5/pgxpool"
)

// PGXConnect establishes a connection to the PostgreSQL database using pgx.
func PGXConnect(dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	log.Println("Successfully connected to the database")
	return conn, nil
}

// func PGXPoolConnect(dsn string) (*pgxpool.Pool, error) {
// 	poolConfig, err := pgxpool.ParseConfig(dsn)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse pgx pool config: %w", err)
// 	}
// 	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to pgx pool: %w", err)
// 	}
// 	log.Panicln("Successfully connected to pgx pool")
// 	return pool, nil
// }

// PGXExec executes a query without returning any rows.
func PGXExec(db *pgx.Conn, ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	if db == nil {
		return pgconn.CommandTag{}, errors.New("database connection is not initialized")
	}
	log.Printf("Executing query: %s with args: %v", query, args)
	tag, err := db.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return pgconn.CommandTag{}, err
	}
	log.Printf("Query executed successfully: %s", query)
	return tag, nil
}

// func PGXPoolExec(pool *pgxpool.Pool, ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
// 	if pool == nil {
// 		return pgconn.CommandTag{}, errors.New("pgx pool is not initialized")
// 	}
// 	log.Printf("Executing pgx query: %s with args: %v", query, args)
// 	tag, err := pool.Exec(ctx, query, args...)
// 	if err != nil {
// 		return pgconn.CommandTag{}, fmt.Errorf("failed to execute pgx query: %w", err)
// 	}
// 	log.Println("pgx query executed successfully")
// 	return tag, nil
// }

// PGXExecNoResult executes a query without returning any rows and without returning the result tag.
func PGXExecNoResult(db *pgx.Conn, ctx context.Context, query string, args ...interface{}) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}
	_, err := PGXExec(db, ctx, query, args...)
	return err
}

// PGXQuery executes a query that returns rows.
func PGXQuery(db *pgx.Conn, ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	log.Printf("Executing query: %s with args: %v", query, args)
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	log.Printf("Query executed successfully: %s", query)
	return rows, nil
}

// PGXQueryRow executes a query that is expected to return at most one row.
func PGXQueryRow(db *pgx.Conn, ctx context.Context, query string, args ...interface{}) pgx.Row {
	if db == nil {
		return nil
	}
	log.Printf("Executing query row: %s with args: %v", query, args)
	return db.QueryRow(ctx, query, args...)
}
