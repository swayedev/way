package database

import (
	"fmt"
)

// CheckDriver checks and returns the appropriate driver
func CheckDriver(driver string) string {
	switch driver {
	case "postgres", "cockroachdb", "pgx":
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

// DriverImportHint returns the optional import package for a normalized SQL driver.
func DriverImportHint(driver string) string {
	switch CheckDriver(driver) {
	case "mysql":
		return `import _ "github.com/swayedev/way/database/drivers/mysql"`
	case "pgx":
		return `import _ "github.com/swayedev/way/database/drivers/pgx"`
	case "sqlite3":
		return `import _ "github.com/swayedev/way/database/drivers/sqlite"`
	case "sqlserver":
		return `import _ "github.com/swayedev/way/database/drivers/sqlserver"`
	case "godror":
		return `import _ "github.com/swayedev/way/database/drivers/godror"`
	case "clickhouse", "firebirdsql":
		return "register a compatible database/sql driver before opening the connection"
	default:
		return "use database.CheckDriver to normalize a supported driver name"
	}
}

// CheckDSN constructs the DSN based on the driver and provided parameters
func CheckDSN(driver, dsn, dbName, dbHost, dbPort, dbUser, dbPass string) string {
	if dsn != "" {
		return dsn
	}
	return setDBDSN(driver, dbName, dbHost, dbPort, dbUser, dbPass)
}

// setDBDSN constructs the DSN based on the driver and provided parameters
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
