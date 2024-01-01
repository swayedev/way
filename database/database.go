package database

import (
	"github.com/swayedev/way/database/config"
	wayPgx "github.com/swayedev/way/database/pgx"
	waySql "github.com/swayedev/way/database/sql"
)

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
