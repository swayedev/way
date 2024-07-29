package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB struct {
	Driver string
	UsePgx bool
	sql    *sql.DB
	pgx    *pgx.Conn
}

// All databases will use the sql.DB interface unless pgx is specified
// pgx is a PostgreSQL driver that is more efficient than the standard sql.DB
// pgx will be used if the driver is set to "pgx" or "postgres" and pgxNew is called
func New() DB {
	return DB{
		Driver: "",
		UsePgx: false,
		sql:    nil,
		pgx:    nil,
	}
}

func (d *DB) SetDB(db interface{}, driver string) {
	switch v := db.(type) {
	case *sql.DB:
		d.SqlNew(v, driver)
	case *pgx.Conn:
		d.PGXNew(v)
	}
}

func (d *DB) PGXNew(db *pgx.Conn) {
	d.pgx = db
	d.UsePgx = true
	d.Driver = "pgx"
}

func (d *DB) SqlNew(db *sql.DB, driver string) {
	d.sql = db
	d.UsePgx = false
	d.Driver = SetDBDriver(driver, false)
}

func (d *DB) Sql() *sql.DB {
	return d.sql
}

func (d *DB) Pgx() *pgx.Conn {
	return d.pgx
}

func (d *DB) Open(dsn string) error {
	if d.UsePgx {
		return d.PgxOpen(dsn)
	}
	return d.SqlOpen(d.Driver, dsn)
}

func (d *DB) PgxOpen(dsn string) error {
	d.UsePgx = true
	d.Driver = "pgx"
	db, err := PGXConnect(dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	d.pgx = db
	return nil
}

func (d *DB) SqlOpen(driver, dsn string) error {
	d.UsePgx = false
	d.Driver = SetDBDriver(driver, d.UsePgx)
	db, err := SQLConnect(d.Driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	d.sql = db
	return nil
}

func (d *DB) Close() error {
	if d.UsePgx {
		return d.PgxClose()
	}

	return d.SqlClose()
}

func (d *DB) PgxClose() error {
	if d.pgx == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return d.pgx.Close(context.Background())
}

func (d *DB) SqlClose() error {
	if d.sql == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return d.sql.Close()
}

func (d *DB) Exec(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	if d.UsePgx {
		return d.PgxExec(ctx, query, args...)
	}

	return d.SqlExec(ctx, query, args...)
}

func (d *DB) PgxExec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return PGXExec(d.pgx, ctx, query, args...)
}

func (d *DB) SqlExec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return SQLExec(d.sql, ctx, query, args...)
}

func (d *DB) ExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	if d.UsePgx {
		return d.PgxExecNoResult(ctx, query, args...)
	}

	return d.SqlExecNoResult(ctx, query, args...)
}

func (d *DB) PgxExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	return PGXExecNoResult(d.pgx, ctx, query, args...)
}

func (d *DB) SqlExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	return SQLExecNoResult(d.sql, ctx, query, args...)
}

func (d *DB) Query(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	if d.UsePgx {
		return d.PgxQuery(ctx, query, args...)
	}
	return d.SqlQuery(ctx, query, args...)
}

func (d *DB) PgxQuery(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return PGXQuery(d.pgx, ctx, query, args...)
}

func (d *DB) SqlQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return SQLQuery(d.sql, ctx, query, args...)
}

func (d *DB) QueryRow(ctx context.Context, query string, args ...interface{}) interface{} {
	if d.UsePgx {
		return d.PgxQueryRow(ctx, query, args...)
	}
	return d.SqlQueryRow(ctx, query, args...)
}

func (d *DB) PgxQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return PGXQueryRow(d.pgx, ctx, query, args...)
}

func (d *DB) SqlQueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return SQLQueryRow(d.sql, ctx, query, args...)
}

func (d *DB) SetDriver(driver string, usePgx bool) {
	d.Driver = SetDBDriver(driver, d.UsePgx)
}

func (d *DB) SetDSN(driver, dsn, dbName, dbHost, dbPort, dbUser, dbPass string) string {
	return CheckDSN(driver, dsn, dbName, dbHost, dbPort, dbUser, dbPass)
}

func (d *DB) CheckDriver() {
	d.Driver = CheckDriver(d.Driver)
}

func Connect(driver, dsn string) (interface{}, error) {
	if driver == "pgx" {
		return PGXConnect(dsn)
	}
	return SQLConnect(driver, dsn)
}

// checks
func CheckDriver(driver string) string {
	switch driver {
	case "pgx", "postgres", "mysql", "sqlite3":
		return driver
	}
	return ""
}

func SetDBDriver(driver string, usePgx bool) string {
	switch driver {
	case "postgres", "cockroachdb":
		if usePgx {
			return "pgx"
		}
		return "postgres"
	case "pgx":
		return "pgx"
	case "mysql", "sql":
		return "mysql"
	case "sqlite", "sqlite3":
		return "sqlite3"
	case "clickhouse", "ch":
		return "clickhouse"
	case "firebirdsql", "firebird":
		return "firebirdsql"
	case "godror", "oracle":
		return "godror"
	case "sqlserver", "mssql":
		return "sqlserver"
	default:
		return ""
	}
}

func CheckDSN(driver, dsn, dbName, dbHost, dbPort, dbUser, dbPass string) string {
	if dsn != "" {
		return dsn
	}
	return setDBDSN(driver, dbName, dbHost, dbPort, dbUser, dbPass)
}

func setDBDSN(driver, dbName, dbHost, dbPort, dbUser, dbPass string) string {
	switch driver {
	case "postgres", "pgx":
		return setPostgresDSN(dbUser, dbPass, dbHost, dbPort, dbName)
	case "mysql":
		return setMysqlDSN(dbUser, dbPass, dbHost, dbPort, dbName)
	case "sqlite3":
		return setSqliteDSN(dbName, dbHost)
	case "clickhouse":
		return setClickhouseDSN(dbUser, dbPass, dbHost, dbPort, dbName)
	case "firebirdsql":
		return setFirebirdDSN(dbUser, dbPass, dbHost, dbPort, dbName)
	case "godror":
		return setOracleDSN(dbUser, dbPass, dbHost, dbPort, dbName)
	case "sqlserver":
		return setSQLServerDSN(dbUser, dbPass, dbHost, dbPort, dbName)
	default:
		return ""
	}
}

func setPostgresDSN(dbName, dbHost, dbPort, dbUser, dbPass string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
}

func setMysqlDSN(dbName, dbHost, dbPort, dbUser, dbPass string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
}

func setClickhouseDSN(dbUser, dbPass, dbHost, dbPort, dbName string) string {
	return fmt.Sprintf("tcp://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
}

func setFirebirdDSN(dbUser, dbPass, dbHost, dbPort, dbName string) string {
	return fmt.Sprintf("%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
}

func setOracleDSN(dbUser, dbPass, dbHost, dbPort, dbName string) string {
	return fmt.Sprintf("%s/%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
}

func setSQLServerDSN(dbUser, dbPass, dbHost, dbPort, dbName string) string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", dbUser, dbPass, dbHost, dbPort, dbName)
}

func setSqliteDSN(dbName, dbHost string) string {
	// if the driver is sqlite3 the host is the path to the database file
	// the dbName is the name of the database file
	return dbHost + dbName + ".db"
}

func SupportedDrivers() []string {
	return []string{
		"clickhouse", "firebirdsql",
		"godror", "mysql",
		"pgx", "postgres",
		"sqlite3", "sqlserver",
	}
}

func DefaultPorts() map[string]string {
	return map[string]string{
		"clickhouse":  "8123",
		"firebirdsql": "3050",
		"godror":      "1521",
		"mysql":       "3306",
		"pgx":         "5432",
		"postgres":    "5432",
		"sqlite3":     "",
		"sqlserver":   "1433",
	}
}

// TODO - Set up a way to pass in a struct and have it automatically
// create the table and columns
// func processStruct(s interface{}) {
// 	val := reflect.ValueOf(s).Elem()
// 	typ := val.Type()

// 	for i := 0; i < val.NumField(); i++ {
// 		field := val.Field(i)
// 		dbTag := typ.Field(i).Tag.Get("db")

// 		if dbTag != "-" {
// 			// Process fields that are not excluded from JSON
// 			fmt.Printf("DB field: %s, Value: %v\n", dbTag, field.Interface())
// 		}
// 	}
// }
