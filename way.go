package way

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type Way struct {
	// startupMutex is used to synchronize startup operations.
	startupMutex sync.RWMutex

	// db is the database connection.
	db *DB

	// router is the HTTP request router.
	router *mux.Router

	// sessions is the session manager.
	sessions *Session

	// Server is the HTTP server.
	Server *http.Server

	// Listener is the network listener.
	Listener net.Listener
}

// HandlerFunc is a function type that represents a handler for a request.
// It takes a *Context parameter, which provides information about the request
// and allows the handler to generate a response.
type HandlerFunc func(*Context)

// MiddlewareFunc represents a function that takes a HandlerFunc and returns a modified HandlerFunc.
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// New creates a new instance of Way.
// It initializes the sessions and sets session defaults if necessary.
// It returns a pointer to the newly created Way instance.
func New() *Way {
	sessions := NewSession()
	if useDefaultSession() {
		setSessionDefaults(sessions)
	}
	return &Way{
		db:       new(DB),
		router:   mux.NewRouter(),
		Server:   new(http.Server),
		sessions: sessions,
	}
}

// SetDB sets the database connection for the Way object.
// It takes a pointer to a DB object as a parameter and assigns it to the db field of the Way object.
func (w *Way) SetDB(db *DB) {
	w.db = db
}

// SetSession sets the session for the Way.
// It takes a pointer to a Session as a parameter and assigns it to the sessions field of the Way.
func (w *Way) SetSession(s *Session) {
	w.sessions = s
}

// Use adds a middleware to the middleware stack.
func (w *Way) Use(middleware ...MiddlewareFunc) {
	w.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
			ctx := NewContext(wr, r, w.db, w.sessions)

			// Create a chain of middleware handlers
			handler := func(c *Context) {
				next.ServeHTTP(wr, r)
			}

			// Loop through the middleware in reverse order and chain them
			for i := len(middleware) - 1; i >= 0; i-- {
				handler = middleware[i](handler)
			}

			// Call the first middleware with the context
			handler(ctx)
		})
	})
}

// adaptHandler adapts a `HandlerFunc` to `http.HandlerFunc`.
func adaptHandler(db *DB, s *Session, handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r, db, s)
		handler(ctx)
	}
}

// handleFuncWithMethod registers a new route with a matcher for the URL path and the HTTP method.
func (w *Way) handleFuncWithMethod(path string, handler HandlerFunc, method string) {
	w.router.HandleFunc(path, adaptHandler(w.db, w.sessions, handler)).Methods(method)
}

// HandleFunc registers a new route with a matcher for the URL path.
func (w *Way) HandleFunc(path string, handler HandlerFunc) {
	w.router.HandleFunc(path, adaptHandler(w.db, w.sessions, handler))
}

// GET is a shortcut for `HandleFunc(path, handler)` for the "GET" method.
func (w *Way) GET(path string, handler HandlerFunc) {
	w.handleFuncWithMethod(path, handler, "GET")
}

// POST is a shortcut for `HandleFunc(path, handler)` for the "POST" method.
func (w *Way) POST(path string, handler HandlerFunc) {
	w.handleFuncWithMethod(path, handler, "POST")
}

// PUT is a shortcut for `HandleFunc(path, handler)` for the "PUT" method.
func (w *Way) PUT(path string, handler HandlerFunc) {
	w.handleFuncWithMethod(path, handler, "PUT")
}

// DELETE is a shortcut for `HandleFunc(path, handler)` for the "DELETE" method.
func (w *Way) DELETE(path string, handler HandlerFunc) {
	w.handleFuncWithMethod(path, handler, "DELETE")
}

// PATCH is a shortcut for `HandleFunc(path, handler)` for the "PATCH" method.
func (w *Way) PATCH(path string, handler HandlerFunc) {
	w.handleFuncWithMethod(path, handler, "PATCH")
}

// OPTIONS is a shortcut for `HandleFunc(path, handler)` for the "OPTIONS" method.
func (w *Way) OPTIONS(path string, handler HandlerFunc) {
	w.handleFuncWithMethod(path, handler, "OPTIONS")
}

// HEAD is a shortcut for `HandleFunc(path, handler)` for the "HEAD" method.
func (w *Way) HEAD(path string, handler HandlerFunc) {
	w.handleFuncWithMethod(path, handler, "HEAD")
}

// newListener creates a new net.Listener.
func newListener(network, address string) (net.Listener, error) {
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// Start starts the server.
func (w *Way) Start(address string) error {
	var err error
	w.startupMutex.Lock()
	w.Listener, err = newListener("tcp", address)
	if err != nil {
		w.startupMutex.Unlock()
		return err
	}
	w.Server.Handler = w.router
	w.startupMutex.Unlock()
	return w.Server.Serve(w.Listener)
}

// Close immediately stops the server.
// It internally calls `http.Server#Close()`.
func (w *Way) Close() error {
	w.startupMutex.Lock()
	defer w.startupMutex.Unlock()
	return w.Server.Close()
}

// Shutdown stops the server gracefully.
// It internally calls `http.Server#Shutdown()`.
func (w *Way) Shutdown(ctx context.Context) error {
	w.startupMutex.Lock()
	defer w.startupMutex.Unlock()
	return w.Server.Shutdown(ctx)
}

// DB Functions
// Db returns the database instance associated with the Way.
func (w *Way) Db() *DB {
	return w.db
}

// private functions
var envStoreName = "WAY_DEFAULT_STORE_NAME"
var envStoreEncryptionKey = "WAY_DEFAULT_STORE_ENCRYPTION_KEY"

var envCookieName = "WAY_DEFAULT_COOKIE_NAME"
var envCookieEncryptionKey = "WAY_DEFAULT_COOKIE_ENCRYPTION_KEY"
var envCookieAuthenticationKey = "WAY_DEFAULT_COOKIE_AUTHENTICATION_KEY"

// useDefaultSession checks if the default session should be used.
// It returns true if either the environment store encryption key is set,
// or both the environment cookie encryption key and the environment cookie authentication key are set.
// Otherwise, it returns false.
func useDefaultSession() bool {
	if envStoreEncryptionKey != "" {
		return true
	}
	if os.Getenv(envCookieEncryptionKey) != "" && os.Getenv(envCookieAuthenticationKey) != "" {
		return true
	}
	return false
}

// setSessionDefaults sets the default values for a Session object.
// It assigns the default store name and default cookie name to the Session object.
// If the environment variable envStoreEncryptionKey is not empty, it creates a new CookieStore with the store encryption key and assigns it to the default store.
// If the environment variables envCookieEncryptionKey and envCookieAuthenticationKey are not empty, it creates a new securecookie with the encryption key and authentication key and assigns it to the default cookie.
func setSessionDefaults(s *Session) {
	s.defaultStore = getDefaultStoreName()
	s.defaultCookie = getDefaultCookieName()
	if envStoreEncryptionKey != "" {
		s.stores[s.defaultStore] = sessions.NewCookieStore(getStoreEncryptionKey())
	}
	if os.Getenv(envCookieEncryptionKey) != "" && os.Getenv(envCookieAuthenticationKey) != "" {
		s.cookies[s.defaultCookie] = securecookie.New(
			getEncryptionKey(),     // Encryption key
			getAuthenticationKey(), // Authentication key
		)
	}
}

// getDefaultCookieName returns the default cookie name.
// If the environment variable envCookieName is not set, it uses the default name 'way'.
func getDefaultCookieName() string {
	if os.Getenv(envCookieName) == "" {
		log.Printf("%s not set, using default name of 'way'", envCookieName)
		return "way"
	}
	return os.Getenv(envCookieName)
}

// getDefaultStoreName returns the default store name.
// If the environment variable 'envStoreName' is not set,
// it logs a message and returns the default name 'way'.
// Otherwise, it returns the value of the environment variable 'envStoreName'.
func getDefaultStoreName() string {
	if os.Getenv(envStoreName) == "" {
		log.Printf("%s not set, using default name of 'way'", envStoreName)
		return "way"
	}
	return os.Getenv(envStoreName)
}

// getEncryptionKey retrieves the encryption key used for cookie encryption.
// It checks the environment variable `envCookieEncryptionKey` and returns its value as a byte slice.
// If the environment variable is not set, it logs a fatal error and terminates the program.
func getEncryptionKey() []byte {
	// check env, else create new and update env
	if os.Getenv(envCookieEncryptionKey) == "" {
		log.Fatalf("%s is required", envCookieEncryptionKey)
	}
	return []byte(os.Getenv(envCookieEncryptionKey))
}

// getAuthenticationKey retrieves the authentication key from the environment variable.
// If the environment variable is not set, it logs a fatal error.
// Returns the authentication key as a byte slice.
func getAuthenticationKey() []byte {
	// check env, else create new and update env
	if os.Getenv(envCookieAuthenticationKey) == "" {
		log.Fatalf("%s is required", envCookieAuthenticationKey)
	}
	return []byte(os.Getenv(envCookieAuthenticationKey))
}

// getStoreEncryptionKey retrieves the encryption key used for storing data in the application's store.
// It checks the environment variable `envStoreEncryptionKey` and returns the key as a byte slice.
// If the environment variable is not set, it logs a fatal error and terminates the program.
func getStoreEncryptionKey() []byte {
	// check env, else create new and update env
	if os.Getenv(envStoreEncryptionKey) == "" {
		log.Fatalf("%s is required", envStoreEncryptionKey)
	}
	return []byte(envStoreEncryptionKey)
}
