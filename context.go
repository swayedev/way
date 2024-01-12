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
	"github.com/swayedev/way/crypto"
)

// Response is the standard go HTTP response writer.
// Request is the standard go HTTP request.
// db is the way database connection.
// Session is the way session.
type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	db       *DB
	Session  *Session
}

func NewContext(w http.ResponseWriter, r *http.Request, d *DB, s *Session) *Context {
	return &Context{Response: w, Request: r, db: d, Session: s}
}

func (c *Context) SetSession(s *Session) {
	c.Session = s
}

func (c *Context) GetDB() *DB {
	return c.db
}

func (c *Context) GetSession(name string) sessions.Store {
	return c.Session.stores[name]
}

func (c *Context) Parms() map[string]string {
	return mux.Vars(c.Request)
}

func (c *Context) SqlExec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return c.db.SqlExec(ctx, query, args...)
}

func (c *Context) SqlExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	return c.db.SqlExecNoResult(ctx, query, args...)
}

func (c *Context) SqlQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return c.db.SqlQuery(ctx, query, args...)
}

func (c *Context) SqlQueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return c.db.SqlQueryRow(ctx, query, args...)
}

func (c *Context) PgxExec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return c.db.PgxExec(ctx, query, args...)
}

func (c *Context) PgxExecNoResult(ctx context.Context, query string, args ...interface{}) error {
	return c.db.PgxExecNoResult(ctx, query, args...)
}

func (c *Context) PgxQuery(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return c.db.PgxQuery(ctx, query, args...)
}

func (c *Context) PgxQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return c.db.PgxQueryRow(ctx, query, args...)
}

func (c *Context) Redirect(code int, url string) {
	http.Redirect(c.Response, c.Request, url, code)
}

func (c *Context) JSON(code int, i interface{}) {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(code)
	if err := json.NewEncoder(c.Response).Encode(i); err != nil {
		// handle error, e.g., log it or send an internal server error
		log.Printf("Error encoding JSON response: %v \n", err)
	}
}

func (c *Context) HTML(code int, htmlContent string) {
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response.WriteHeader(code)
	if _, err := c.Response.Write([]byte(htmlContent)); err != nil {
		log.Printf("Error writing HTML response: %v \n", err)
	}
}

func (c *Context) String(code int, i interface{}) {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(code)
	switch v := i.(type) {
	case string:
		c.Response.Write([]byte(v))
	case []byte:
		c.Response.Write(v)
	default:
		log.Printf("Error encoding String response: %v \n", v)
	}
}

func (c *Context) XML(code int, i interface{}) {
	c.Response.Header().Set("Content-Type", "application/xml")
	c.Response.WriteHeader(code)
	if err := xml.NewEncoder(c.Response).Encode(i); err != nil {
		log.Printf("Error encoding XML response: %v \n", err)
	}

}

func (c *Context) Data(code int, data []byte) {
	c.Response.WriteHeader(code)
	if _, err := c.Response.Write(data); err != nil {
		log.Printf("Error writing data: %v \n", err)
	}
}

func (c *Context) Status(code int) {
	c.Response.WriteHeader(code)
}

func (c *Context) Image(code int, contentType string, imageData []byte) {
	c.Response.Header().Set("Content-Type", contentType)
	c.Response.WriteHeader(code)
	if _, err := c.Response.Write(imageData); err != nil {
		log.Printf("Error writing image data: %v \n", err)
	}
}

func (c *Context) SetHeader(key string, value string) {
	c.Response.Header().Set(key, value)
}

func (c *Context) ProxyMedia(mediaURL string) {
	resp, err := http.Get(mediaURL)
	if err != nil {
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
	io.Copy(c.Response, resp.Body)
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Response, cookie)
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
