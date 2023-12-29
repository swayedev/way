package way

import (
	"database/sql"
	"net/http"
)

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	Sql      *sql.DB
}

func NewContext(s *sql.DB, w http.ResponseWriter, r *http.Request) *Context {
	return &Context{Response: w, Request: r, Sql: s}
}
