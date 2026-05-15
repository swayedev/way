module github.com/swayedev/way

go 1.26.0

require (
	github.com/go-sql-driver/mysql v1.10.0
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/securecookie v1.1.2
	github.com/gorilla/sessions v1.4.0
	github.com/jackc/pgx/v5 v5.9.2
	github.com/mattn/go-sqlite3 v1.14.44
	github.com/swayedev/fcrypt v1.0.0-rc1
	golang.org/x/crypto v0.51.0
)

replace github.com/swayedev/fcrypt => ../fcrypt

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/VictoriaMetrics/easyproto v1.2.0 // indirect
	github.com/go-logfmt/logfmt v0.6.1 // indirect
	github.com/godror/knownpb v0.3.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/planetscale/vtprotobuf v0.6.0 // indirect
	golang.org/x/exp v0.0.0-20260508232706-74f9aab9d74a // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

require (
	github.com/denisenkom/go-mssqldb v0.12.3
	github.com/godror/godror v0.50.0
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/text v0.37.0 // indirect
)
