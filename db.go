package way

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	wayPgx "github.com/swayedev/way/database/pgx"
	waySql "github.com/swayedev/way/database/sql"
)

type DB struct {
	Driver string
	sql    *sql.DB
	pgx    *pgx.Conn
}

// Sql Functions
func (d *DB) Sql() *sql.DB {
	return d.sql
}

func (d *DB) SqlNew(db *sql.DB) {
	d.sql = db
	d.Driver = "mysql"
}

func (d *DB) SqlOpen() error {
	db, err := waySql.Connect()
	if err != nil {
		log.Println("Failed to open database connection:", err)
		return err
	}

	d.sql = db
	d.Driver = "mysql"
	return nil
}

func (d *DB) SqlClose() error {
	if d.sql != nil {
		d.Driver = ""
		return d.sql.Close()
	}
	return nil
}

func (d *DB) SqlExec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return sqlExec(d.sql, ctx, query, args...)
}

func (d *DB) SqlExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	return sqlExecNoResult(d.sql, ctx, query, args...)
}

func (d *DB) SqlQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return sqlQuery(d.sql, ctx, query, args...)
}

func (d *DB) SqlQueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return sqlQueryRow(d.sql, ctx, query, args...)
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

// Pgx Functions
func (d *DB) Pgx() *pgx.Conn {
	return d.pgx
}

func (d *DB) PgxNew(db *pgx.Conn) {
	d.pgx = db
	d.Driver = "postgres"
}

func (d *DB) PgxOpen() error {
	db, err := wayPgx.Connect()
	if err != nil {
		log.Println("Failed to open database connection:", err)
		return err
	}

	d.pgx = db
	d.Driver = "postgres"
	return nil
}

func (d *DB) PgxClose() error {
	if d.sql != nil {
		d.Driver = ""
		return d.pgx.Close(context.Background())
	}
	return nil
}

func (d *DB) PgxExec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return pgxExec(d.pgx, ctx, query, args...)
}

func (d *DB) PgxExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	return pgxExecNoResult(d.pgx, ctx, query, args...)
}

func (d *DB) PgxQuery(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return pgxQuery(d.pgx, ctx, query, args...)
}

func (d *DB) PgxQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return pgxQueryRow(d.pgx, ctx, query, args...)
}

func pgxExec(db *pgx.Conn, ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	if db == nil {
		return pgconn.CommandTag{}, errors.New("database connection is not initialized")
	}
	return db.Exec(ctx, query, args...)
}

func pgxExecNoResult(db *pgx.Conn, ctx context.Context, query string, args ...interface{}) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}
	_, err := pgxExec(db, ctx, query, args...)
	return err
}

func pgxQuery(db *pgx.Conn, ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	return db.Query(ctx, query, args...)
}

func pgxQueryRow(db *pgx.Conn, ctx context.Context, query string, args ...interface{}) pgx.Row {
	if db == nil {
		return nil
	}
	return db.QueryRow(ctx, query, args...)
}
