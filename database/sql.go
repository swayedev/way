package database

import (
	"context"
	"database/sql"
	"errors"
	"log"

	_ "github.com/denisenkom/go-mssqldb" // Microsoft SQL Server
	_ "github.com/go-sql-driver/mysql"   // MySQL
	_ "github.com/godror/godror"         // Oracle
	_ "github.com/jackc/pgx/v5/stdlib"   // PostgreSQL
	_ "github.com/mattn/go-sqlite3"      // SQLite
)

// SQLConnect opens a connection to the specified database and checks if it is alive.
func SQLConnect(driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		return nil, err
	}

	log.Printf("Successfully connected to the database with driver %s", driver)
	return db, nil
}

// SQLExec executes a query without returning any rows.
func SQLExec(db *sql.DB, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	log.Printf("Executing query: %s with args: %v", query, args)
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	log.Printf("Query executed successfully: %s", query)
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
	log.Printf("Executing query: %s with args: %v", query, args)
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
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
