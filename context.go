package way

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/swayedev/way/crypto"
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
func NewContext(w http.ResponseWriter, r *http.Request, d *DB, s *Session, l *log.Logger, cr *Crypto) *Context {
	return &Context{Response: w, Request: r, db: d, Session: s, Logger: l, Crypto: cr}
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

// MultipartForm returns the parsed multipart form.
func (c *Context) MultipartForm() (*multipart.Form, error) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		return nil, err
	}
	return c.Request.MultipartForm, nil
}

// Bind binds the request body into provided type `i`.
func (c *Context) Bind(i interface{}) error {
	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "" {
		return errors.New("Content-Type header is missing")
	}
	switch {
	case contentType == "application/json":
		return json.NewDecoder(c.Request.Body).Decode(i)
	case contentType == "application/xml" || contentType == "text/xml":
		return xml.NewDecoder(c.Request.Body).Decode(i)
	case contentType == "application/x-www-form-urlencoded":
		if err := c.Request.ParseForm(); err != nil {
			return err
		}
		decoder := schema.NewDecoder()
		return decoder.Decode(i, c.Request.PostForm)
	default:
		return errors.New("Unsupported Content-Type")
	}
}

// QueryParam returns the query param for the provided name.
func (c *Context) QueryParam(name string) string {
	return c.Request.URL.Query().Get(name)
}

// QueryParams returns the query parameters as `url.Values`.
func (c *Context) QueryParams() url.Values {
	return c.Request.URL.Query()
}

// QueryString returns the URL query string.
func (c *Context) QueryString() string {
	return c.Request.URL.RawQuery
}

// FormValue returns the form field value for the provided name.
func (c *Context) FormValue(name string) string {
	return c.Request.FormValue(name)
}

// FormParams returns the form parameters as `url.Values`.
func (c *Context) FormParams() (url.Values, error) {
	if err := c.Request.ParseForm(); err != nil {
		return nil, err
	}
	return c.Request.PostForm, nil
}

// FormFile returns the multipart form file for the provided name.
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	_, fileHeader, err := c.Request.FormFile(name)
	return fileHeader, err
}

// Cookies returns the HTTP cookies sent with the request.
func (c *Context) Cookies() []*http.Cookie {
	return c.Request.Cookies()
}

// Request Handling Functions

// IsTLS returns true if HTTP connection is TLS otherwise false.
func (c *Context) IsTLS() bool {
	return c.Request.TLS != nil
}

// IsWebSocket returns true if HTTP connection is WebSocket otherwise false.
func (c *Context) IsWebSocket() bool {
	return c.Request.Header.Get("Upgrade") == "websocket"
}

// Scheme returns the HTTP protocol scheme, `http` or `https`.
func (c *Context) Scheme() string {
	if c.Request.TLS != nil {
		return "https"
	}
	return "http"
}

// RealIP returns the client's network address based on `X-Forwarded-For` or `X-Real-IP` request header.
func (c *Context) RealIP() string {
	xff := c.Request.Header.Get("X-Forwarded-For")
	if xff != "" {
		return xff
	}
	xri := c.Request.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}
	return c.Request.RemoteAddr
}

// Path returns the registered path for the handler.
func (c *Context) Path() string {
	return c.Request.URL.Path
}

// SetPath sets the registered path for the handler.
func (c *Context) SetPath(p string) {
	c.Request.URL.Path = p
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

// Cryptographic Functions (Crypto) - Not Implemented
func (c *Context) SetCrypto(cr *Crypto) {
	c.Crypto = cr
}

// Deprecated: Use the `crypto` package interface instead.
// HashStringToString hashes a string to another string.
func (c *Context) HashStringToString(value string) string {
	return crypto.HashStringToString(value)
}

// Deprecated: Use the `crypto` package interface instead.
// // HashString hashes a string to a byte array.
func (c *Context) HashString(value string) [32]byte {
	return crypto.HashString(value)
}

// Deprecated: Use the `crypto` package interface instead.
// // HashByte hashes a byte array to another byte array.
func (c *Context) HashByte(value []byte) [32]byte {
	return crypto.HashByte(value)
}

// Deprecated: Use the `crypto` package interface instead.
// // Encrypt encrypts data using a passphrase.
func (c *Context) Encrypt(data []byte, passphrase string) (string, error) {
	encrypted, err := crypto.Encrypt(data, passphrase)
	if err != nil {
		c.Logger.Printf("Error encrypting data: %v", err)
	}
	return encrypted, err
}

// Deprecated: Use the `crypto` package interface instead.
// // Decrypt decrypts data using a passphrase.
func (c *Context) Decrypt(encrypted string, passphrase string) ([]byte, error) {
	decrypted, err := crypto.Decrypt(encrypted, passphrase)
	if err != nil {
		c.Logger.Printf("Error decrypting data: %v", err)
	}
	return decrypted, err
}
