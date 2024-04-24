package sql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"github.com/swayedev/way/database/config"
)

// Connect to database
func Connect(driver string, uri string) (*sql.DB, error) {
	if driver == "" && config.GetDbType() == "" {
		return nil, fmt.Errorf("database driver is not set")
	}

	// driver overrides config
	if driver == "" {
		driver = config.GetDbType()
	}

	if uri == "" && config.GetDbUri() == "" {
		switch driver {
		case "postgres":
			uri = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.GetDbUser(), config.GetDbPassword(), config.GetDbHost(), config.GetDbPort(), config.GetDbName())
		case "mysql":
			uri = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.GetDbUser(), config.GetDbPassword(), config.GetDbHost(), config.GetDbPort(), config.GetDbName())
		case "sqlite3":
			uri = fmt.Sprintf("%s", config.GetDbName())
		}
	}

	if uri == "" {
		uri = config.GetDbUri()
	}

	sqlDriver := driver
	if driver == "postgres" {
		sqlDriver = "pgx"
	}

	// Open database connection
	db, err := sql.Open(sqlDriver, uri)
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
