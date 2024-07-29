package database

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/denisenkom/go-mssqldb" // Microsoft SQL Server
	_ "github.com/go-sql-driver/mysql"   // MySQL
	_ "github.com/godror/godror"         // Oracle
	_ "github.com/jackc/pgx/v5/stdlib"   // PostgreSQL
	_ "github.com/mattn/go-sqlite3"      // SQLite
)

// Connect to database
func SQLConnect(driver, dsn string) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	// Check if database is alive
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SQLExec(db *sql.DB, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	return db.ExecContext(ctx, query, args...)
}

func SQLExecNoResult(db *sql.DB, ctx context.Context, query string, args ...interface{}) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}
	_, err := SQLExec(db, ctx, query, args...)
	return err
}

func SQLQuery(db *sql.DB, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	return db.QueryContext(ctx, query, args...)
}

func SQLQueryRow(db *sql.DB, ctx context.Context, query string, args ...interface{}) *sql.Row {
	if db == nil {
		return nil
	}
	return db.QueryRowContext(ctx, query, args...)
}
