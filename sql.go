package way

import (
	"context"
	"database/sql"
	"errors"
	"log"

	_ "github.com/go-sql-driver/mysql"
	waySql "github.com/swayedev/way/database/sql"
)

func (w *Way) SQL() *sql.DB {
	return w.sql
}

func (w *Way) SqlNew(db *sql.DB) {
	w.sql = db
}

func (w *Way) SqlOpen() error {
	db, err := waySql.Connect()
	if err != nil {
		log.Println("Failed to open database connection:", err)
		return err
	}

	w.sql = db
	return nil
}

func (w *Way) SqlClose() error {
	if w.sql != nil {
		return w.sql.Close()
	}
	return nil
}

func (w *Way) SqlExec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return sqlExec(w.sql, ctx, query, args...)
}

func (w *Way) SqlExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	return sqlExecNoResult(w.sql, ctx, query, args...)
}

func (w *Way) SqlQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return sqlQuery(w.sql, ctx, query, args...)
}

func (w *Way) SqlQueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return sqlQueryRow(w.sql, ctx, query, args...)
}

func sqlExec(db *sql.DB, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	return db.ExecContext(ctx, query, args...)
}

func sqlExecNoResult(db *sql.DB, ctx context.Context, query string, args ...interface{}) error {
	_, err := sqlExec(db, ctx, query, args...)
	return err
}

func sqlQuery(db *sql.DB, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	return db.QueryContext(ctx, query, args...)
}

func sqlQueryRow(db *sql.DB, ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.QueryRowContext(ctx, query, args...)
}
