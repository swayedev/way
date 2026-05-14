package way

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/swayedev/way/crypto"
)

// Response is the standard go HTTP response writer.
// Request is the standard go HTTP request.
// db is the way database connection.
// Session is the way session.
type Context struct {
	Response   http.ResponseWriter
	Request    *http.Request
	db         *DB
	Session    *Session
	Logger     *log.Logger
	HTTPClient *http.Client
}

func NewContext(w http.ResponseWriter, r *http.Request, d *DB, s *Session, l *log.Logger) *Context {
	return newContextWithHTTPClient(w, r, d, s, l, defaultHTTPClient())
}

func newContextWithHTTPClient(w http.ResponseWriter, r *http.Request, d *DB, s *Session, l *log.Logger, client *http.Client) *Context {
	if l == nil {
		l = defaultLogger()
	}
	if client == nil {
		client = defaultHTTPClient()
	}
	return &Context{Response: w, Request: r, db: d, Session: s, Logger: l, HTTPClient: client}
}
func (c *Context) Log() *log.Logger {
	if c.Logger != nil {
		return c.Logger
	}
	return defaultLogger()
}
func (c *Context) SetSession(s *Session) {
	c.Session = s
}

func (c *Context) GetDB() *DB {
	return c.db
}

func (c *Context) GetSession(name string) sessions.Store {
	store, err := c.GetSessionE(name)
	if err != nil {
		c.Log().Printf("Session store not found: %s", name)
		return nil
	}
	c.Log().Printf("Session store retrieved: %s", name)
	return store
}

func (c *Context) GetSessionE(name string) (sessions.Store, error) {
	if c.Session == nil {
		return nil, errors.New("session manager is not initialized")
	}
	return c.Session.StoreE(name)
}

// Parms returns a map of string parameters associated with the http request context.
// The keys of the map are the parameter names, and the values are the parameter values.
func (c *Context) Parms() map[string]string {
	params := mux.Vars(c.Request)
	return params
}

// Parm returns the value of the specified parameter from the http request context.
// I.E. in the codebase - way.GET("http://example.com/user/{userId}",...)
// via a browser or curl request - http://example.com/user/1234
// c.Parm("userId") returns "1234".
// In the above example, the parameter name is "userId" and the parameter value is "1234".
// If the parameter does not exist, an empty string is returned.
func (c *Context) Parm(param string) string {
	return c.Parms()[param]
}

func (c *Context) SqlExec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	c.Log().Printf("Executing SQL query")
	result, err := c.db.SQLExec(ctx, query, args...)
	if err != nil {
		c.Log().Printf("Error executing SQL query: %v", err)
	}
	return result, err
}

func (c *Context) SqlExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	c.Log().Printf("Executing SQL query with no result")
	err := c.db.SQLExecNoResult(ctx, query, args...)
	if err != nil {
		c.Log().Printf("Error executing SQL query with no result: %v", err)
	}
	return err
}

func (c *Context) SqlQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	c.Log().Printf("Executing SQL query")
	rows, err := c.db.SQLQuery(ctx, query, args...)
	if err != nil {
		c.Log().Printf("Error executing SQL query: %v", err)
	}
	return rows, err
}

func (c *Context) SqlQueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	c.Log().Printf("Executing SQL query row")
	row := c.db.SQLQueryRow(ctx, query, args...)
	return row
}

func (c *Context) PgxExec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	c.Log().Printf("Executing PGX query")
	commandTag, err := c.db.PGXExec(ctx, query, args...)
	if err != nil {
		c.Log().Printf("Error executing PGX query: %v", err)
	}
	return commandTag, err
}

func (c *Context) PgxExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	c.Log().Printf("Executing PGX query with no result")
	err := c.db.PGXExecNoResult(ctx, query, args...)
	if err != nil {
		c.Log().Printf("Error executing PGX query with no result: %v", err)
	}
	return err
}

func (c *Context) PgxQuery(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	c.Log().Printf("Executing PGX query")
	rows, err := c.db.PGXQuery(ctx, query, args...)
	if err != nil {
		c.Log().Printf("Error executing PGX query: %v", err)
	}
	return rows, err
}

func (c *Context) PgxQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	c.Log().Printf("Executing PGX query row")
	row := c.db.PGXQueryRow(ctx, query, args...)
	return row
}

func (c *Context) Redirect(code int, url string) {
	c.Log().Printf("Redirecting with status code: %d", code)
	http.Redirect(c.Response, c.Request, url, code)
}

func (c *Context) JSON(code int, i interface{}) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(i); err != nil {
		c.Log().Printf("Error encoding JSON response: %v", err)
		http.Error(c.Response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(code)
	if _, err := c.Response.Write(body.Bytes()); err != nil {
		c.Log().Printf("Error writing JSON response: %v", err)
	}
}

// HTML sends an HTML response with the specified status code.
func (c *Context) HTML(code int, htmlContent string) {
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response.WriteHeader(code)
	if _, err := c.Response.Write([]byte(htmlContent)); err != nil {
		c.Log().Printf("Error writing HTML response: %v", err)
	}
}

func (c *Context) String(code int, i interface{}) {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(code)
	switch v := i.(type) {
	case string:
		_, err := c.Response.Write([]byte(v))
		if err != nil {
			c.Log().Printf("Error writing string response: %v", err)
		}
	case []byte:
		_, err := c.Response.Write(v)
		if err != nil {
			c.Log().Printf("Error writing byte response: %v", err)
		}
	default:
		c.Log().Printf("Error encoding string response: unsupported type %T", v)
	}
}

func (c *Context) XML(code int, i interface{}) {
	var body bytes.Buffer
	if err := xml.NewEncoder(&body).Encode(i); err != nil {
		c.Log().Printf("Error encoding XML response: %v", err)
		http.Error(c.Response, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	c.Response.Header().Set("Content-Type", "application/xml")
	c.Response.WriteHeader(code)
	if _, err := c.Response.Write(body.Bytes()); err != nil {
		c.Log().Printf("Error writing XML response: %v", err)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Response.WriteHeader(code)
	if _, err := c.Response.Write(data); err != nil {
		c.Log().Printf("Error writing data: %v", err)
	}
}

func (c *Context) Status(code int) {
	c.Response.WriteHeader(code)
	c.Log().Printf("Status set to %d", code)
}

func (c *Context) Image(code int, contentType string, imageData []byte) {
	c.Response.Header().Set("Content-Type", contentType)
	c.Response.WriteHeader(code)
	if _, err := c.Response.Write(imageData); err != nil {
		c.Log().Printf("Error writing image data: %v", err)
	}
}

func (c *Context) SetHeader(key string, value string) {
	c.Response.Header().Set(key, value)
	c.Log().Printf("Header set: %s", key)
}

func (c *Context) ProxyMedia(mediaURL string) {
	client := c.HTTPClient
	if client == nil {
		client = defaultHTTPClient()
	}
	resp, err := client.Get(mediaURL)
	if err != nil {
		c.Log().Printf("Failed to fetch media: %v", err)
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
		c.Log().Printf("Error streaming media: %v", err)
	}
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Response, cookie)
	c.Log().Printf("Cookie set: %s", cookie.Name)
}

func (c *Context) GetCookie(name string) (*http.Cookie, error) {
	return c.Request.Cookie(name)
}

func (c *Context) DeleteCookie(name string) {
	c.SetCookie(&http.Cookie{
		Name:   name,
		MaxAge: -1,
	})
}

// func (c *Context) GetSession(name string) (*http.Cookie, error) {
// 	return c.Request.Cookie(name)
// }

// func (c *Context) SetSession(cookie *http.Cookie) {
// 	http.SetCookie(c.Response, cookie)
// }

// func (c *Context) DeleteSession(name string) {
// 	c.SetSession(&http.Cookie{
// 		Name:   name,
// 		MaxAge: -1,
// 	})
// }

// func (c *Context) GetSessionValue(name string) (string, error) {
// 	cookie, err := c.GetSession(name)
// 	if err != nil {
// 		return "", err
// 	}
// 	return cookie.Value, nil
// }

// func (c *Context) SetSessionValue(name string, value string) {
// 	c.SetSession(&http.Cookie{
// 		Name:  name,
// 		Value: value,
// 	})
// }

// func (c *Context) DeleteSessionValue(name string) {
// 	c.DeleteSession(name)
// }

func (c *Context) HashStringToString(value string) string {
	return crypto.HashStringToString(value)
}

func (c *Context) HashString(value string) [32]byte {
	return crypto.HashString(value)
}

func (c *Context) HashByte(value []byte) [32]byte {
	return crypto.HashByte(value)
}

func (c *Context) Encrypt(data []byte, passphrase string) (string, error) {
	return crypto.Encrypt(data, passphrase)
}

func (c *Context) Decrypt(encrypted string, passphrase string) ([]byte, error) {
	return crypto.Decrypt(encrypted, passphrase)
}

var (
	SqlErrNoRows = sql.ErrNoRows
	PgxErrNoRows = pgx.ErrNoRows
)
