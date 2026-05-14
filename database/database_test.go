package database

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	_ "github.com/swayedev/way/database/drivers/sqlite"
)

func TestCheckDriverNormalizesSupportedDrivers(t *testing.T) {
	tests := map[string]string{
		"postgres": "pgx",
		"pgx":      "pgx",
		"mysql":    "mysql",
		"sqlite":   "sqlite3",
		"sqlite3":  "sqlite3",
		"mssql":    "sqlserver",
		"oracle":   "godror",
		"unknown":  "",
	}

	for input, want := range tests {
		if got := CheckDriver(input); got != want {
			t.Fatalf("CheckDriver(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestDriverImportHintMentionsOptionalPackage(t *testing.T) {
	hint := DriverImportHint("sqlite")
	if !strings.Contains(hint, "github.com/swayedev/way/database/drivers/sqlite") {
		t.Fatalf("hint = %q, want sqlite optional driver package", hint)
	}
}

func TestSQLConnectUnsupportedDriverReturnsHelpfulError(t *testing.T) {
	_, err := SQLConnect("unsupported-driver", "unused")
	if err == nil {
		t.Fatal("SQLConnect() error = nil, want unsupported driver error")
	}
	if !strings.Contains(err.Error(), "use database.CheckDriver") {
		t.Fatalf("error = %q, want driver hint", err.Error())
	}
}

func TestSQLiteSQLHelpersWithOptionalDriver(t *testing.T) {
	db, err := SQLConnect("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("SQLConnect() error = %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	if _, err := SQLExec(db, ctx, "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)"); err != nil {
		t.Fatalf("SQLExec(create) error = %v", err)
	}
	if _, err := SQLExec(db, ctx, "INSERT INTO users(name) VALUES (?)", "Ada"); err != nil {
		t.Fatalf("SQLExec(insert) error = %v", err)
	}
	row := SQLQueryRow(db, ctx, "SELECT name FROM users WHERE id = ?", 1)
	if row == nil {
		t.Fatal("SQLQueryRow() = nil")
	}
	var name string
	if err := row.Scan(&name); err != nil && err != sql.ErrNoRows {
		t.Fatalf("Scan() error = %v", err)
	}
	if name != "Ada" {
		t.Fatalf("name = %q, want Ada", name)
	}
}
