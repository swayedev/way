package sql

import (
	"database/sql"
	"fmt"

	"github.com/swayedev/way/database/config"
)

// Connect to database
func Connect() (*sql.DB, error) {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.GetDbUser(), config.GetDbPassword(), config.GetDbHost(), config.GetDbPort(), config.GetDbName())
	// Open database connection
	db, err := sql.Open("mysql", uri)
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
