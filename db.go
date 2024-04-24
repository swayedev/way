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

func (d *DB) New(db interface{}) {
	switch v := db.(type) {
	case *sql.DB:
		d.SqlNew(v, "sql")
	case *pgx.Conn:
		d.PgxNew(v)
	}
}

func (d *DB) Open() error {
	switch d.Driver {
	case "pgx":
		return d.PgxOpen()
	case "postgres":
		fallthrough
	case "sql":
		fallthrough
	case "mysql":
		fallthrough
	case "sqlite3":
		return d.SqlOpen(d.Driver, "")
	}
	return nil
}

func (d *DB) Close() error {
	switch d.Driver {
	case "pgx":
		return d.PgxClose()
	case "postgres":
		fallthrough
	case "mysql":
		fallthrough
	case "sqlite3":
		return d.SqlClose()
	}
	return nil
}

func (d *DB) Exec(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	switch d.Driver {
	case "pgx":
		return d.PgxExec(ctx, query, args...)
	case "postgres":
		fallthrough
	case "mysql":
		fallthrough
	case "sqlite3":
		return d.SqlExec(ctx, query, args...)
	}
	return nil, errors.New("database driver is not initialized")
}

func (d *DB) ExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	switch d.Driver {
	case "pgx":
		return d.PgxExecNoResult(ctx, query, args...)
	case "postgres":
		fallthrough
	case "mysql":
		fallthrough
	case "sqlite3":
		return d.SqlExecNoResult(ctx, query, args...)
	}
	return errors.New("database driver is not initialized")
}

func (d *DB) Query(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	switch d.Driver {
	case "pgx":
		return d.PgxQuery(ctx, query, args...)
	case "postgres":
		fallthrough
	case "mysql":
		fallthrough
	case "sqlite3":
		return d.SqlQuery(ctx, query, args...)
	}
	return nil, errors.New("database driver is not initialized")
}

func (d *DB) QueryRow(ctx context.Context, query string, args ...interface{}) interface{} {
	switch d.Driver {
	case "pgx":
		return d.PgxQueryRow(ctx, query, args...)
	case "postgres":
		fallthrough
	case "mysql":
		fallthrough
	case "sqlite3":
		return d.SqlQueryRow(ctx, query, args...)
	}
	return nil
}

// Sql Functions
func (d *DB) Sql() *sql.DB {
	return d.sql
}

func (d *DB) SqlNew(db *sql.DB, driver string) {
	d.sql = db
	// mysql, sqlite3, postgres
	d.Driver = driver
}

func (d *DB) SqlOpen(driver string, uri string) error {
	d.Driver = driver

	db, err := waySql.Connect(driver, uri)
	if err != nil {
		log.Println("Failed to open database connection:", err)
		return err
	}

	d.sql = db
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
	if db == nil {
		return errors.New("database connection is not initialized")
	}
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
	if db == nil {
		return nil
	}
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
