package database

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb" // Microsoft SQL Server
	_ "github.com/go-sql-driver/mysql"   // MySQL
	_ "github.com/godror/godror"         // Oracle
	_ "github.com/jackc/pgx/v5/stdlib"   // PostgreSQL
	_ "github.com/mattn/go-sqlite3"      // SQLite
)

// DriverConfig holds the configuration for different database drivers.
type DriverConfig struct {
	Driver string
	DSN    string
	DBName string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
}

type DBConfig interface {
	// SetDSN sets the DSN for the database connection.
	GetDriver() string
	GetDSN() string
	GetUsePooling() bool
	GetMaxOpenConns() int
	GetMaxIdleConns() int
	GetConnMaxLifetime() time.Duration
}

// DBError represents a custom error type for database operations.
type DBError struct {
	Op  string // Operation
	Err error  // Original error
}

// Error implements the error interface for DBError.
func (e *DBError) Error() string {
	return fmt.Sprintf("db error during %s: %v", e.Op, e.Err)
}

// Unwrap returns the underlying error for DBError.
func (e *DBError) Unwrap() error {
	return e.Err
}

// DBError operation constants
const (
	OpConfigParse = "parse config"
	OpDriverCheck = "check driver"
	OpDSNSet      = "set DSN"
	// Database operations
	OpPing    = "ping"
	OpConnect = "connect"
	OpExec    = "execute"
	OpQuery   = "query"
	OpClose   = "close"
	// Pooling operations
	OpPoolConnect = "connect with pool"
	// Transaction operations
	OpTxBegin    = "begin transaction"
	OpTxCommit   = "commit transaction"
	OpTxRollback = "rollback transaction"
)

// DBError error messages
const (
	ErrUnsupportedDriver      = "unsupported database driver"
	ErrDatabaseNotInitialized = "database connection is not initialized"
)

// NewDBError creates a new DBError.
func NewDBError(op string, err error) error {
	return &DBError{Op: op, Err: err}
}

// CheckDriver checks and returns the appropriate driver, or an error if unsupported
func CheckDriver(driver string) (string, error) {
	switch driver {
	case "postgres", "cockroachdb", "pgx":
		return "pgx", nil
	case "mysql", "sql":
		return "mysql", nil
	case "sqlite", "sqlite3":
		return "sqlite3", nil
	case "clickhouse", "ch":
		return "clickhouse", nil
	case "firebirdsql", "firebird":
		return "firebirdsql", nil
	case "godror", "oracle":
		return "godror", nil
	case "sqlserver", "mssql":
		return "sqlserver", nil
	default:
		return "", NewDBError(OpDriverCheck, errors.New(ErrUnsupportedDriver))
	}
}

// CheckDSN constructs the DSN based on the driver and provided parameters, or returns an error.
func CheckDSN(config DriverConfig) (string, error) {
	if config.DSN != "" {
		return config.DSN, nil
	}
	return setDBDSN(config)
}

// setDBDSN constructs the DSN based on the driver and provided parameters, or returns an error.
func setDBDSN(config DriverConfig) (string, error) {
	switch config.Driver {
	case "postgres", "pgx":
		return setPostgresDSN(config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName), nil
	case "mysql":
		return setMysqlDSN(config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName), nil
	case "sqlite3":
		return setSqliteDSN(config.DBName, config.DBHost), nil
	case "clickhouse":
		return setClickhouseDSN(config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName), nil
	case "firebirdsql":
		return setFirebirdDSN(config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName), nil
	case "godror":
		return setOracleDSN(config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName), nil
	case "sqlserver":
		return setSQLServerDSN(config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName), nil
	default:
		return "", NewDBError(OpDSNSet, errors.New(ErrUnsupportedDriver))
	}
}

// setPostgresDSN constructs the DSN for PostgreSQL
func setPostgresDSN(dbUser, dbPass, dbHost, dbPort, dbName string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
}

// setMysqlDSN constructs the DSN for MySQL
func setMysqlDSN(dbUser, dbPass, dbHost, dbPort, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
}

// setClickhouseDSN constructs the DSN for ClickHouse
func setClickhouseDSN(dbUser, dbPass, dbHost, dbPort, dbName string) string {
	return fmt.Sprintf("tcp://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
}

// setFirebirdDSN constructs the DSN for FirebirdSQL
func setFirebirdDSN(dbUser, dbPass, dbHost, dbPort, dbName string) string {
	return fmt.Sprintf("%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
}

// setOracleDSN constructs the DSN for Oracle
func setOracleDSN(dbUser, dbPass, dbHost, dbPort, dbName string) string {
	return fmt.Sprintf("%s/%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
}

// setSQLServerDSN constructs the DSN for SQL Server
func setSQLServerDSN(dbUser, dbPass, dbHost, dbPort, dbName string) string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", dbUser, dbPass, dbHost, dbPort, dbName)
}

// setSqliteDSN constructs the DSN for SQLite
func setSqliteDSN(dbName, dbHost string) string {
	// if the driver is sqlite3 the host is the path to the database file
	// the dbName is the name of the database file
	return dbHost + dbName + ".db"
}

// SupportedDrivers returns a list of supported database drivers
func SupportedDrivers() []string {
	return []string{
		"clickhouse", "firebirdsql",
		"godror", "mysql", "pgx",
		"sqlite3", "sqlserver",
	}
}

// DefaultPorts returns a map of default ports for supported database drivers
func DefaultPorts() map[string]string {
	return map[string]string{
		"clickhouse":  "8123",
		"firebirdsql": "3050",
		"godror":      "1521",
		"mysql":       "3306",
		"pgx":         "5432",
		"cockroachdb": "26257",
		"sqlite3":     "",
		"sqlserver":   "1433",
	}
}
