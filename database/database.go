package database

import (
	"database/sql"

	"github.com/swayedev/way/database/config"
	waySql "github.com/swayedev/way/database/sql"
)

func Connect() (*sql.DB, error) {
	switch config.GetDbType() {
	case "mysql":
		return waySql.Connect()
	}
	return waySql.Connect()
}
