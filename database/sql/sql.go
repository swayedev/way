package sql

import (
	"database/sql"
	"fmt"

	"github.com/swayedev/way/database/config"
)

// Connect to database
func Connect(t string, uri string) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open(t, uri)
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

// Connect to a MySql database
func MySqlConnect() (*sql.DB, error) {
	uri := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.GetDbUser(),
		config.GetDbPassword(),
		config.GetDbHost(),
		config.GetDbPort(),
		config.GetDbName(),
	)

	return Connect("mysql", uri)
}
