package way

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Context represents the context for a request.
type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	db       *DB
	Session  *Session
	Logger   *log.Logger
	Crypto   *Crypto
}

// NewContext creates a new Context instance.
func NewContext(w http.ResponseWriter, r *http.Request, d *DB, s *Session, l *log.Logger) *Context {
	return &Context{Response: w, Request: r, db: d, Session: s, Logger: l}
}

// Log returns the logger.
func (c *Context) Log() *log.Logger {
	return c.Logger
}

// SetSession sets the session.
func (c *Context) SetSession(s *Session) {
	c.Session = s
}

// GetDB returns the database instance.
func (c *Context) GetDB() *DB {
	return c.db
}

// GetSession returns a session store by name.
func (c *Context) GetSession(name string) sessions.Store {
	store := c.Session.stores[name]
	if store == nil {
		c.Logger.Printf("Session store not found: %s", name)
	} else {
		c.Logger.Printf("Session store retrieved: %s", name)
	}
	return store
}

// Parms returns the request parameters.
func (c *Context) Parms() map[string]string {
	params := mux.Vars(c.Request)
	c.Logger.Printf("Request parameters: %v", params)
	return params
}

// SQL Execution Functions

// SqlExec executes an SQL query and returns the result.
func (c *Context) SqlExec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	c.Logger.Printf("Executing SQL query: %s with args: %v", query, args)
	result, err := c.db.SQLExec(ctx, query, args...)
	if err != nil {
		c.Logger.Printf("Error executing SQL query: %v", err)
	}
	return result, err
}

// SqlExecNoResult executes an SQL query without returning a result.
func (c *Context) SqlExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	c.Logger.Printf("Executing SQL query with no result: %s with args: %v", query, args)
	err := c.db.SQLExecNoResult(ctx, query, args...)
	if err != nil {
		c.Logger.Printf("Error executing SQL query with no result: %v", err)
	}
	return err
}

// SqlQuery executes an SQL query and returns rows.
func (c *Context) SqlQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	c.Logger.Printf("Executing SQL query: %s with args: %v", query, args)
	rows, err := c.db.SQLQuery(ctx, query, args...)
	if err != nil {
		c.Logger.Printf("Error executing SQL query: %v", err)
	}
	return rows, err
}

// SqlQueryRow executes an SQL query that is expected to return at most one row.
func (c *Context) SqlQueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	c.Logger.Printf("Executing SQL query row: %s with args: %v", query, args)
	row := c.db.SQLQueryRow(ctx, query, args...)
	return row
}

// PGX Execution Functions

// PgxExec executes a PGX query and returns the command tag.
func (c *Context) PgxExec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	c.Logger.Printf("Executing PGX query: %s with args: %v", query, args)
	commandTag, err := c.db.PGXExec(ctx, query, args...)
	if err != nil {
		c.Logger.Printf("Error executing PGX query: %v", err)
	}
	return commandTag, err
}

// PgxExecNoResult executes a PGX query without returning a result.
func (c *Context) PgxExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	c.Logger.Printf("Executing PGX query with no result: %s with args: %v", query, args)
	err := c.db.PGXExecNoResult(ctx, query, args...)
	if err != nil {
		c.Logger.Printf("Error executing PGX query with no result: %v", err)
	}
	return err
}

// PgxQuery executes a PGX query and returns rows.
func (c *Context) PgxQuery(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	c.Logger.Printf("Executing PGX query: %s with args: %v", query, args)
	rows, err := c.db.PGXQuery(ctx, query, args...)
	if err != nil {
		c.Logger.Printf("Error executing PGX query: %v", err)
	}
	return rows, err
}

// PgxQueryRow executes a PGX query that is expected to return at most one row.
func (c *Context) PgxQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	c.Logger.Printf("Executing PGX query row: %s with args: %v", query, args)
	row := c.db.PGXQueryRow(ctx, query, args...)
	return row
}

// Response Functions

// Redirect redirects the request to a provided URL with status code.
func (c *Context) Redirect(code int, url string) {
	c.Logger.Printf("Redirecting to URL: %s with status code: %d", url, code)
	http.Redirect(c.Response, c.Request, url, code)
}

// JSON sends a JSON response with status code.
func (c *Context) JSON(code int, i interface{}) {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(code)
	if err := json.NewEncoder(c.Response).Encode(i); err != nil {
		c.Logger.Printf("Error encoding JSON response: %v", err)
		http.Error(c.Response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// HTML sends an HTML response with the specified status code.
func (c *Context) HTML(code int, htmlContent string) {
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response.WriteHeader(code)
	if _, err := c.Response.Write([]byte(htmlContent)); err != nil {
		c.Logger.Printf("Error writing HTML response: %v", err)
		http.Error(c.Response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// String sends a string response with status code.
func (c *Context) String(code int, i interface{}) {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(code)
	switch v := i.(type) {
	case string:
		_, err := c.Response.Write([]byte(v))
		if err != nil {
			c.Logger.Printf("Error writing string response: %v", err)
		}
	case []byte:
		_, err := c.Response.Write(v)
		if err != nil {
			c.Logger.Printf("Error writing byte response: %v", err)
		}
	default:
		c.Logger.Printf("Error encoding string response: unsupported type %v", v)
	}
}

// XML sends an XML response with status code.
func (c *Context) XML(code int, i interface{}) {
	c.Response.Header().Set("Content-Type", "application/xml")
	c.Response.WriteHeader(code)
	if err := xml.NewEncoder(c.Response).Encode(i); err != nil {
		c.Logger.Printf("Error encoding XML response: %v", err)
		http.Error(c.Response, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Data sends a response with raw data and status code.
func (c *Context) Data(code int, data []byte) {
	c.Response.WriteHeader(code)
	if _, err := c.Response.Write(data); err != nil {
		c.Logger.Printf("Error writing data: %v", err)
	}
}

// Status sets the status code for the response.
func (c *Context) Status(code int) {
	c.Response.WriteHeader(code)
	c.Logger.Printf("Status set to %d", code)
}

// Image sends an image response with status code and content type.
func (c *Context) Image(code int, contentType string, imageData []byte) {
	c.Response.Header().Set("Content-Type", contentType)
	c.Response.WriteHeader(code)
	if _, err := c.Response.Write(imageData); err != nil {
		c.Logger.Printf("Error writing image data: %v", err)
	}
}

// SetHeader sets a header key-value pair.
func (c *Context) SetHeader(key, value string) {
	c.Response.Header().Set(key, value)
	c.Logger.Printf("Header set: %s = %s", key, value)
}

// ProxyMedia proxies media from a given URL.
func (c *Context) ProxyMedia(mediaURL string) {
	resp, err := http.Get(mediaURL)
	if err != nil {
		c.Logger.Printf("Failed to fetch media from URL: %s, error: %v", mediaURL, err)
		http.Error(c.Response, "Failed to fetch media", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Response.Header().Add(key, value)
		}
	}

	// Stream the content
	c.Response.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(c.Response, resp.Body); err != nil {
		c.Logger.Printf("Error streaming media: %v", err)
	}
}

// SetCookie sets a cookie.
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Response, cookie)
	c.Logger.Printf("Cookie set: %s", cookie.String())
}

// GetCookie retrieves a cookie by name.
func (c *Context) GetCookie(name string) (*http.Cookie, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		c.Logger.Printf("Error retrieving cookie: %v", err)
	}
	return cookie, err
}

// DeleteCookie deletes a cookie by name.
func (c *Context) DeleteCookie(name string) {
	c.SetCookie(&http.Cookie{
		Name:   name,
		MaxAge: -1,
	})
}
