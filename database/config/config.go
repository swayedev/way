package config

import (
	"log"
	"os"
)

var ports = map[string]string{
	"mysql":    "3306",
	"postgres": "5432",
}

func GetDbType() string {
	if dbType := os.Getenv("WAY_DB_TYPE"); dbType != "" {
		return dbType
	}
	log.Println("Environment variable WAY_DB_TYPE is not set")
	return "mysql"
}

func GetDbUser() string {
	if user := os.Getenv("WAY_DB_USER"); user != "" {
		return user
	}
	log.Println("Environment variable WAY_DB_USER is not set")
	return ""
}

func GetDbPassword() string {
	if password := os.Getenv("WAY_DB_PASSWORD"); password != "" {
		return password
	}
	log.Println("Environment variable WAY_DB_PASSWORD is not set")
	return ""
}

func GetDbHost() string {
	if host := os.Getenv("WAY_DB_HOST"); host != "" {
		return host
	}
	log.Println("Environment variable WAY_DB_HOST is not set")
	return ""
}

func GetDbPort() string {
	if port := os.Getenv("WAY_DB_PORT"); port != "" {
		return port
	}
	log.Println("Environment variable WAY_DB_PORT is not set")
	return ports[GetDbType()]
}

func GetDbName() string {
	if name := os.Getenv("WAY_DB_NAME"); name != "" {
		return name
	}
	log.Println("Environment variable WAY_DB_NAME is not set")
	return ""
}
