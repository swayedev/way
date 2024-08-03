package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PGXConfig holds the configuration for PostgreSQL connection and pooling.
type PGXConfig struct {
	DSN             string
	UsePooling      bool
	MaxConns        int32
	MinConns        int32
	ConnMaxLifetime time.Duration
}

// PGXConnect establishes a connection to the PostgreSQL database using pgx.
func PGXConnect(config PGXConfig) (*pgx.Conn, *pgxpool.Pool, error) {
	if config.UsePooling {
		return connectWithPool(config)
	}
	return connectWithoutPool(config)
}

// connectWithoutPool establishes a direct connection to the PostgreSQL database.
func connectWithoutPool(config PGXConfig) (*pgx.Conn, *pgxpool.Pool, error) {
	conn, err := pgx.Connect(context.Background(), config.DSN)
	if err != nil {
		return nil, nil, NewDBError(OpConnect, err)
	}
	log.Println("Successfully connected to the database")
	return conn, nil, nil
}

// connectWithPool establishes a connection pool to the PostgreSQL database.
func connectWithPool(config PGXConfig) (*pgx.Conn, *pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, nil, NewDBError(OpConfigParse, err)
	}

	poolConfig.MaxConns = config.MaxConns
	poolConfig.MinConns = config.MinConns
	poolConfig.MaxConnLifetime = config.ConnMaxLifetime

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, nil, NewDBError(OpPoolConnet, err)
	}
	log.Println("Successfully connected to the database with pooling")
	return nil, pool, nil
}

// PGXExec executes a query without returning any rows.
func PGXExec(db *pgx.Conn, pool *pgxpool.Pool, ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	if db == nil && pool == nil {
		return pgconn.CommandTag{}, NewDBError(OpExec, errors.New(ErrDatabaseNotInitialized))
	}

	var err error
	var tag pgconn.CommandTag

	if db != nil {
		tag, err = db.Exec(ctx, query, args...)
	} else {
		tag, err = pool.Exec(ctx, query, args...)
	}

	if err != nil {
		return pgconn.CommandTag{}, NewDBError(OpExec, err)
	}

	log.Printf("Query executed successfully: %s", query)
	return tag, nil
}

// PGXExecNoResult executes a query without returning any rows and without returning the result tag.
func PGXExecNoResult(db *pgx.Conn, pool *pgxpool.Pool, ctx context.Context, query string, args ...interface{}) error {
	_, err := PGXExec(db, pool, ctx, query, args...)
	return err
}

// PGXQuery executes a query that returns rows.
func PGXQuery(db *pgx.Conn, pool *pgxpool.Pool, ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	if db == nil && pool == nil {
		return nil, NewDBError(OpQuery, errors.New(ErrDatabaseNotInitialized))
	}

	var err error
	var rows pgx.Rows

	if db != nil {
		rows, err = db.Query(ctx, query, args...)
	} else {
		rows, err = pool.Query(ctx, query, args...)
	}

	if err != nil {
		return nil, NewDBError(OpQuery, err)
	}

	log.Printf("Query executed successfully: %s", query)
	return rows, nil
}

// PGXQueryRow executes a query that is expected to return at most one row.
func PGXQueryRow(db *pgx.Conn, pool *pgxpool.Pool, ctx context.Context, query string, args ...interface{}) pgx.Row {
	if db == nil && pool == nil {
		return nil
	}

	if db != nil {
		return db.QueryRow(ctx, query, args...)
	}

	return pool.QueryRow(ctx, query, args...)
}

// PGXBeginTx begins a new transaction.
func PGXBeginTx(db *pgx.Conn, pool *pgxpool.Pool, ctx context.Context) (pgx.Tx, error) {
	if db == nil && pool == nil {
		return nil, NewDBError(OpTxBegin, errors.New(ErrDatabaseNotInitialized))
	}

	var err error
	var tx pgx.Tx

	if db != nil {
		tx, err = db.Begin(ctx)
	} else {
		tx, err = pool.Begin(ctx)
	}

	if err != nil {
		return nil, NewDBError(OpTxBegin, err)
	}

	return tx, nil
}

// PGXCommitTx commits the given transaction.
func PGXCommitTx(tx pgx.Tx) error {
	if err := tx.Commit(context.Background()); err != nil {
		return NewDBError(OpTxCommit, err)
	}
	return nil
}

// PGXRollbackTx rolls back the given transaction.
func PGXRollbackTx(tx pgx.Tx) error {
	if err := tx.Rollback(context.Background()); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		return NewDBError(OpTxRollback, err)
	}
	return nil
}
