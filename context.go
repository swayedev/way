package way

import (
	"context"
	"database/sql"
	"net/http"
)

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	sql      *sql.DB
}

func NewContext(s *sql.DB, w http.ResponseWriter, r *http.Request) *Context {
	return &Context{Response: w, Request: r, sql: s}
}

func (c *Context) SqlExec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return sqlExec(c.sql, ctx, query, args...)
}

func (c *Context) SqlExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	return sqlExecNoResult(c.sql, ctx, query, args...)
}

func (c *Context) SqlQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return sqlQuery(c.sql, ctx, query, args...)
}

func (c *Context) SqlQueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return sqlQueryRow(c.sql, ctx, query, args...)
}

func (c *Context) Redirect(url string, code int) {
	http.Redirect(c.Response, c.Request, url, code)
}

func (c *Context) JSON(code int, i interface{}) {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(code)
	c.Response.Write(i.([]byte))
}

func (c *Context) HTML(code int, i interface{}) {
	c.Response.Header().Set("Content-Type", "text/html")
	c.Response.WriteHeader(code)
	c.Response.Write(i.([]byte))
}

func (c *Context) Text(code int, i interface{}) {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(code)
	c.Response.Write(i.([]byte))
}

func (c *Context) XML(code int, i interface{}) {
	c.Response.Header().Set("Content-Type", "application/xml")
	c.Response.WriteHeader(code)
	c.Response.Write(i.([]byte))
}

func (c *Context) Data(code int, i interface{}) {
	c.Response.WriteHeader(code)
	c.Response.Write(i.([]byte))
}

func (c *Context) Status(code int) {
	c.Response.WriteHeader(code)
}

func (c *Context) Image(code int, i interface{}) {
	c.Response.Header().Set("Content-Type", "image/*")
	c.Response.WriteHeader(code)
	c.Response.Write(i.([]byte))
}

func (c *Context) SetHeader(key string, value string) {
	c.Response.Header().Set(key, value)
}
