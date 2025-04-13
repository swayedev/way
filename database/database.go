package database

import (
	"database/sql"

	"github.com/swayedev/way/database/config"
	wayPgx "github.com/swayedev/way/database/pgx"
	waySql "github.com/swayedev/way/database/sql"
)

type Rows interface {
	Delete() bool
	Insert() bool
	RowsAffected() int64
	Select() bool
	String() string
	Update() bool
	Close() error
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
	Err() error
	Next() bool
	NextResultSet() bool
	Scan(dest ...any) error
}

func Connect() (interface{}, error) {
	switch config.GetDbType() {
	case "postgres":
		return wayPgx.Connect()
	}
	return waySql.Connect()
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
