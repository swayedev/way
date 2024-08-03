package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb" // Microsoft SQL Server
	_ "github.com/go-sql-driver/mysql"   // MySQL
	_ "github.com/godror/godror"         // Oracle
	_ "github.com/jackc/pgx/v5/stdlib"   // PostgreSQL
	_ "github.com/mattn/go-sqlite3"      // SQLite
)

// SQLConfig holds the configuration for SQL connection and pooling.
type SQLConfig struct {
	Driver          string
	DSN             string
	UsePooling      bool
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (c SQLConfig) GetDriver() string {
	return c.Driver
}

func (c SQLConfig) GetDSN() string {
	return c.DSN
}

func (c SQLConfig) GetUsePooling() bool {
	return c.UsePooling
}

func (c SQLConfig) GetMaxOpenConns() int {
	return c.MaxOpenConns
}

func (c SQLConfig) GetMaxIdleConns() int {
	return c.MaxIdleConns
}

func (c SQLConfig) GetConnMaxLifetime() time.Duration {
	return c.ConnMaxLifetime
}

// SQLConnect establishes a connection to the SQL database.
func SQLConnect(config DBConfig) (*sql.DB, error) {
	db, err := sql.Open(config.GetDriver(), config.GetDSN())
	if err != nil {
		return nil, NewDBError(OpConnect, err)
	}

	if err = db.Ping(); err != nil {
		return nil, NewDBError(OpPing, err)
	}

	if config.GetUsePooling() {
		db.SetMaxOpenConns(config.GetMaxOpenConns())
		db.SetMaxIdleConns(config.GetMaxIdleConns())
		db.SetConnMaxLifetime(config.GetConnMaxLifetime())
		log.Printf("Database connection pool configured with maxOpenConns=%d, maxIdleConns=%d, connMaxLifetime=%s",
			config.GetMaxOpenConns(), config.GetMaxIdleConns(), config.GetConnMaxLifetime())
	} else {
		log.Println("Database connected without pooling")
	}

	return db, nil
}

// SQLExec executes a query without returning any rows.
func SQLExec(db *sql.DB, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if db == nil {
		return nil, NewDBError(OpExec, errors.New(ErrDatabaseNotInitialized))
	}
	log.Printf("Executing query: %s with args: %v", query, args)
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, NewDBError(OpExec, err)
	}
	log.Printf("Query executed successfully: %s", query)
	return result, nil
}

// SQLExecNoResult executes a query without returning any rows and without returning the result.
func SQLExecNoResult(db *sql.DB, ctx context.Context, query string, args ...interface{}) error {
	_, err := SQLExec(db, ctx, query, args...)
	return err
}

// SQLQuery executes a query that returns rows.
func SQLQuery(db *sql.DB, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if db == nil {
		return nil, NewDBError(OpQuery, errors.New(ErrDatabaseNotInitialized))
	}
	log.Printf("Executing query: %s with args: %v", query, args)
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, NewDBError(OpQuery, err)
	}
	log.Printf("Query executed successfully: %s", query)
	return rows, nil
}

// SQLQueryRow executes a query that is expected to return at most one row.
func SQLQueryRow(db *sql.DB, ctx context.Context, query string, args ...interface{}) *sql.Row {
	if db == nil {
		return nil
	}
	log.Printf("Executing query row: %s with args: %v", query, args)
	return db.QueryRowContext(ctx, query, args...)
}

// SQLBeginTx begins a new transaction.
func SQLBeginTx(db *sql.DB, ctx context.Context) (*sql.Tx, error) {
	if db == nil {
		return nil, NewDBError(OpTxBegin, errors.New(ErrDatabaseNotInitialized))
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, NewDBError(OpTxBegin, err)
	}
	return tx, nil
}

// SQLCommitTx commits the given transaction.
func SQLCommitTx(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return NewDBError(OpTxCommit, err)
	}
	return nil
}

// SQLRollbackTx rolls back the given transaction.
func SQLRollbackTx(tx *sql.Tx) error {
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		return NewDBError(OpTxRollback, err)
	}
	return nil
}
