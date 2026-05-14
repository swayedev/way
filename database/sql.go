package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// SQLConnect opens a connection to the specified database and checks if it is alive.
func SQLConnect(driver, dsn string) (*sql.DB, error) {
	if driver == "" {
		return nil, errors.New("database driver is not set")
	}
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return nil, fmt.Errorf("open %s database connection: %w; %s", driver, err, DriverImportHint(driver))
	}

	if err = db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		db.Close()
		return nil, fmt.Errorf("ping %s database connection: %w; %s", driver, err, DriverImportHint(driver))
	}

	log.Printf("Successfully connected to the database with driver %s", driver)
	return db, nil
}

// SQLExec executes a query without returning any rows.
func SQLExec(db *sql.DB, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	log.Printf("Executing SQL statement")
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	log.Printf("SQL statement executed successfully")
	return result, nil
}

// SQLExecNoResult executes a query without returning any rows and without returning the result.
func SQLExecNoResult(db *sql.DB, ctx context.Context, query string, args ...interface{}) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}
	_, err := SQLExec(db, ctx, query, args...)
	return err
}

// SQLQuery executes a query that returns rows.
func SQLQuery(db *sql.DB, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	log.Printf("Executing SQL query")
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	log.Printf("SQL query executed successfully")
	return rows, nil
}

// SQLQueryRow executes a query that is expected to return at most one row.
func SQLQueryRow(db *sql.DB, ctx context.Context, query string, args ...interface{}) *sql.Row {
	if db == nil {
		return nil
	}
	log.Printf("Executing SQL query row")
	return db.QueryRowContext(ctx, query, args...)
}
