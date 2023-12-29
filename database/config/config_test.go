package config

import (
	"os"
	"testing"
)

func TestGetDbType(t *testing.T) {
	expectedDbType := "postgres"
	os.Setenv("WAY_DB_TYPE", expectedDbType)
	dbType := GetDbType()
	if dbType != expectedDbType {
		t.Errorf("GetDbType() = %v; want %v", dbType, expectedDbType)
	}
}

func TestGetDbUser(t *testing.T) {
	expectedUser := "username"
	os.Setenv("WAY_DB_USER", expectedUser)
	user := GetDbUser()
	if user != expectedUser {
		t.Errorf("GetDbUser() = %v; want %v", user, expectedUser)
	}
}

func TestGetDbPassword(t *testing.T) {
	expectedPassword := "password"
	os.Setenv("WAY_DB_PASSWORD", expectedPassword)
	password := GetDbPassword()
	if password != expectedPassword {
		t.Errorf("GetDbPassword() = %v; want %v", password, expectedPassword)
	}
}

func TestGetDbHost(t *testing.T) {
	expectedHost := "localhost"
	os.Setenv("WAY_DB_HOST", expectedHost)
	host := GetDbHost()
	if host != expectedHost {
		t.Errorf("GetDbHost() = %v; want %v", host, expectedHost)
	}
}

func TestGetDbPort(t *testing.T) {
	expectedPort := "5432"
	os.Setenv("WAY_DB_PORT", expectedPort)
	port := GetDbPort()
	if port != expectedPort {
		t.Errorf("GetDbPort() = %v; want %v", port, expectedPort)
	}
}

func TestGetDbName(t *testing.T) {
	expectedName := "waydb"
	os.Setenv("WAY_DB_NAME", expectedName)
	name := GetDbName()
	if name != expectedName {
		t.Errorf("GetDbName() = %v; want %v", name, expectedName)
	}
}
