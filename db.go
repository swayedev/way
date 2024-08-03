package way

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/swayedev/way/database"
)

// DB struct to handle both sql.DB and pgx.Conn
type DB struct {
	Driver          string
	UsePgx          bool
	sql             *sql.DB
	pgx             *pgx.Conn
	MaxOpenConns    int           // For connection pooling configuration
	MaxIdleConns    int           // For connection pooling configuration
	ConnMaxLifetime time.Duration // For connection pooling configuration
	logger          *log.Logger
}

// NewDB initializes a new DB instance with default connection pool settings.
func NewDB() DB {
	return DB{
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

// NewDBPool initializes a new DB instance with specified connection pool settings.
func NewDBPool(maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration) DB {
	return DB{
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
	}
}

// SetDB sets the database connection based on the type.
func (d *DB) SetDB(db interface{}, driver string) {
	switch v := db.(type) {
	case *sql.DB:
		d.SQLNew(v, driver)
	case *pgx.Conn:
		d.PGXNew(v)
	}
}

// PGXNew initializes a pgx connection.
func (d *DB) PGXNew(db *pgx.Conn) {
	d.pgx = db
	d.UsePgx = true
	d.Driver = "pgx"
}

// SQLNew initializes a sql.DB connection.
func (d *DB) SQLNew(db *sql.DB, driver string) {
	d.sql = db
	d.UsePgx = false
	d.Driver, _ = database.CheckDriver(driver)
}

// SQL returns the sql.DB connection.
func (d *DB) SQL() *sql.DB {
	return d.sql
}

// PGX returns the pgx.Conn connection.
func (d *DB) PGX() *pgx.Conn {
	return d.pgx
}

// Open opens a database connection based on the driver type.
func (d *DB) Open(dsn string) error {
	if d.UsePgx {
		return d.PGXOpen(dsn)
	}
	return d.SQLOpen(d.Driver, dsn)
}

// PGXOpen opens a pgx connection.
func (d *DB) PGXOpen(dsn string) error {
	d.UsePgx = true
	d.Driver = "pgx"
	db, _, err := database.PGXConnect(database.PGXConfig{DSN: dsn, UsePooling: false})
	if err != nil {
		return database.NewDBError(database.OpConnect, fmt.Errorf("failed to open pgx database connection: %w", err))
	}
	d.pgx = db
	return nil
}

// SQLOpen opens a sql.DB connection.
func (d *DB) SQLOpen(driver, dsn string) error {
	d.UsePgx = false
	d.Driver, _ = database.CheckDriver(driver)
	db, err := database.SQLConnect(database.SQLConfig{Driver: driver, DSN: dsn, UsePooling: false})
	if err != nil {
		return database.NewDBError(database.OpConnect, fmt.Errorf("failed to open sql database connection: %w", err))
	}
	db.SetMaxOpenConns(d.MaxOpenConns)
	db.SetMaxIdleConns(d.MaxIdleConns)
	db.SetConnMaxLifetime(d.ConnMaxLifetime)
	d.sql = db
	return nil
}

// Close closes the database connection.
func (d *DB) Close() error {
	if d.UsePgx {
		return d.PGXClose()
	}
	return d.SQLClose()
}

// PGXClose closes the pgx connection.
func (d *DB) PGXClose() error {
	if d.pgx == nil {
		return database.NewDBError(database.OpClose, errors.New("pgx database connection is not initialized"))
	}
	return d.pgx.Close(context.Background())
}

// SQLClose closes the sql.DB connection.
func (d *DB) SQLClose() error {
	if d.sql == nil {
		return database.NewDBError(database.OpClose, errors.New("sql database connection is not initialized"))
	}
	return d.sql.Close()
}

// Exec executes a query.
func (d *DB) Exec(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	if d.UsePgx {
		return d.PGXExec(ctx, query, args...)
	}
	return d.SQLExec(ctx, query, args...)
}

// PGXExec executes a pgx query.
func (d *DB) PGXExec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	if d.pgx == nil {
		return pgconn.CommandTag{}, database.NewDBError(database.OpExec, errors.New("pgx database connection is not initialized"))
	}
	return database.PGXExec(d.pgx, nil, ctx, query, args...)
}

// SQLExec executes a sql.DB query.
func (d *DB) SQLExec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if d.sql == nil {
		return nil, database.NewDBError(database.OpExec, errors.New("sql database connection is not initialized"))
	}
	return database.SQLExec(d.sql, ctx, query, args...)
}

// ExecNoResult executes a query without returning a result.
func (d *DB) ExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	if d.UsePgx {
		return d.PGXExecNoResult(ctx, query, args...)
	}
	return d.SQLExecNoResult(ctx, query, args...)
}

// PGXExecNoResult executes a pgx query without returning a result.
func (d *DB) PGXExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	if d.pgx == nil {
		return database.NewDBError(database.OpExec, errors.New("pgx database connection is not initialized"))
	}
	_, err := database.PGXExec(d.pgx, nil, ctx, query, args...)
	return err
}

// SQLExecNoResult executes a sql.DB query without returning a result.
func (d *DB) SQLExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	if d.sql == nil {
		return database.NewDBError(database.OpExec, errors.New("sql database connection is not initialized"))
	}
	_, err := database.SQLExec(d.sql, ctx, query, args...)
	return err
}

// Query executes a query and returns rows.
func (d *DB) Query(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	if d.UsePgx {
		return d.PGXQuery(ctx, query, args...)
	}
	return d.SQLQuery(ctx, query, args...)
}

// PGXQuery executes a pgx query and returns rows.
func (d *DB) PGXQuery(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	if d.pgx == nil {
		return nil, database.NewDBError(database.OpQuery, errors.New("pgx database connection is not initialized"))
	}
	return database.PGXQuery(d.pgx, nil, ctx, query, args...)
}

// SQLQuery executes a sql.DB query and returns rows.
func (d *DB) SQLQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if d.sql == nil {
		return nil, database.NewDBError(database.OpQuery, errors.New("sql database connection is not initialized"))
	}
	return database.SQLQuery(d.sql, ctx, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
func (d *DB) QueryRow(ctx context.Context, query string, args ...interface{}) interface{} {
	if d.UsePgx {
		return d.PGXQueryRow(ctx, query, args...)
	}
	return d.SQLQueryRow(ctx, query, args...)
}

// PGXQueryRow executes a pgx query that is expected to return at most one row.
func (d *DB) PGXQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	if d.pgx == nil {
		return nil
	}
	return database.PGXQueryRow(d.pgx, nil, ctx, query, args...)
}

// SQLQueryRow executes a sql.DB query that is expected to return at most one row.
func (d *DB) SQLQueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if d.sql == nil {
		return nil
	}
	return database.SQLQueryRow(d.sql, ctx, query, args...)
}

// SetDriver sets the database driver.
func (d *DB) SetDriver(driver string, usePgx bool) {
	if usePgx {
		d.UsePgx = true
		d.Driver = "pgx"
		return
	}

	d.Driver, _ = database.CheckDriver(driver)
}

// SetDSN sets the Data Source Name (DSN).
func (d *DB) SetDSN(driver, dsn, dbName, dbHost, dbPort, dbUser, dbPass string) (string, error) {
	return database.CheckDSN(database.DriverConfig{
		Driver: driver,
		DSN:    dsn,
		DBName: dbName,
		DBHost: dbHost,
		DBPort: dbPort,
		DBUser: dbUser,
		DBPass: dbPass,
	})
}
