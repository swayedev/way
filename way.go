package way

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/swayedev/way/database"

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
	// Logger is the logger.
	Logger *log.Logger
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
		router:   mux.NewRouter(),
		Server:   new(http.Server),
		sessions: sessions,
		Logger:   defaultLogger(),
	}
}

// SetLogger sets the logger.
func (w *Way) SetLogger(logger *log.Logger) {
	w.Logger = logger
}

// SetRouter sets the HTTP router.
func (w *Way) SetRouter(router *mux.Router) {
	w.router = router
}

// SetServer sets the HTTP server.
func (w *Way) SetServer(server *http.Server) {
	w.Server = server
}

// SetListener sets the network listener.
func (w *Way) SetListener(listener net.Listener) {
	w.Listener = listener
}

// Log returns the logger.
func (w *Way) Log() *log.Logger {
	if w.Logger != nil {
		return w.Logger
	}
	return defaultLogger()
}

// SetDB sets the database connection for the Way object.
// It takes a pointer to a DB object as a parameter and assigns it to the db field of the Way object.
func (w *Way) SetDB(db *DB) {
	w.db = db
}

// InitDBFromConfig initializes the database connection from environment variables.
func (w *Way) InitDBFromConfig() error {
	usePGX := GetEnv(envDBUsePGX, "") == "true"
	driver := database.CheckDriver(GetEnv(envDBDriver, ""))
	if driver == "" {
		return errors.New("database driver is not set")
	}
	dsn := database.CheckDSN(driver, GetEnv(envDBDSN, ""), GetEnv(envDBName, ""), GetEnv(envDBHost, ""), GetEnv(envDBPort, ""), GetEnv(envDBUser, ""), GetEnv(envDBPassword, ""))
	if dsn == "" {
		return errors.New("database DSN is not set")
	}
	if usePGX {
		db, err := database.PGXConnect(dsn)
		if err != nil {
			return err
		}
		w.db = &DB{pgx: db, UsePgx: true, Driver: "pgx"}
		return nil
	}
	db, err := database.SQLConnect(driver, dsn)
	if err != nil {
		return err
	}
	w.db = &DB{sql: db, UsePgx: false, Driver: driver}
	return nil
}

// SetDBConnection sets a new database connection with the given driver.
func (w *Way) SetDBConnection(db interface{}, driver string) {
	d := NewDB()
	w.db = &d
	w.db.SetDB(db, driver)
}

// SetSession sets the session.
func (w *Way) SetSession(s *Session) {
	w.sessions = s
}

// Use adds a middleware to the middleware stack.
func (w *Way) Use(middleware ...MiddlewareFunc) {
	w.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
			ctx := NewContext(wr, r, w.db, w.sessions, w.Logger)

			handler := func(c *Context) {
				next.ServeHTTP(wr, r)
			}

			for i := len(middleware) - 1; i >= 0; i-- {
				handler = middleware[i](handler)
			}

			handler(ctx)
		})
	})
}

// adaptHandler adapts a HandlerFunc to http.HandlerFunc.
func adaptHandler(db *DB, s *Session, l *log.Logger, handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r, db, s, l)
		handler(ctx)
	}
}

// handleFuncWithMethod registers a new route with a matcher for the URL path and the HTTP method.
func (w *Way) handleFuncWithMethod(path string, handler HandlerFunc, method string) {
	w.Logger.Printf("Registering route %s", path)
	w.router.HandleFunc(path, adaptHandler(w.db, w.sessions, w.Logger, handler)).Methods(method)
}

// HandleFunc registers a new route with a matcher for the URL path.
func (w *Way) HandleFunc(path string, handler HandlerFunc) {
	w.Logger.Printf("Registering route %s", path)
	w.router.HandleFunc(path, adaptHandler(w.db, w.sessions, w.Logger, handler))
}

// HTTP method shortcuts
func (w *Way) GET(path string, handler HandlerFunc)  { w.handleFuncWithMethod(path, handler, "GET") }
func (w *Way) POST(path string, handler HandlerFunc) { w.handleFuncWithMethod(path, handler, "POST") }
func (w *Way) PUT(path string, handler HandlerFunc)  { w.handleFuncWithMethod(path, handler, "PUT") }
func (w *Way) DELETE(path string, handler HandlerFunc) {
	w.handleFuncWithMethod(path, handler, "DELETE")
}
func (w *Way) PATCH(path string, handler HandlerFunc) { w.handleFuncWithMethod(path, handler, "PATCH") }
func (w *Way) OPTIONS(path string, handler HandlerFunc) {
	w.handleFuncWithMethod(path, handler, "OPTIONS")
}
func (w *Way) HEAD(path string, handler HandlerFunc) { w.handleFuncWithMethod(path, handler, "HEAD") }

// newListener creates a new net.Listener.
func newListener(network, address string) (net.Listener, error) {
	return net.Listen(network, address)
}

// Start starts the server.
func (w *Way) Start(address string) error {
	var err error
	w.startupMutex.Lock()
	defer w.startupMutex.Unlock()

	w.Listener, err = newListener("tcp", address)
	if err != nil {
		return err
	}
	w.Server.Handler = loggingMiddleware(w.Logger, w.router)
	w.Logger.Printf("Server started at %s", address)
	asciiArt := `
	__        ______   __
	\ \      / /  \ \ / /
	 \ \ /\ / / /\ \ V / 
	  \ V  V / /__\ | |  
	   \_/\_/_/----\|_|  
	`
	w.Log().Println(asciiArt)
	return w.Server.Serve(w.Listener)
}

// Close immediately stops the server.
func (w *Way) Close() error {
	w.Logger.Println("Server stopping...")
	w.startupMutex.Lock()
	defer w.startupMutex.Unlock()
	w.Logger.Println("Server stopped")
	return w.Server.Close()
}

// Shutdown stops the server gracefully.
func (w *Way) Shutdown(ctx context.Context) error {
	w.Logger.Println("Server stopping gracefully...")
	w.startupMutex.Lock()
	defer w.startupMutex.Unlock()
	w.Logger.Println("Server stopped gracefully")
	return w.Server.Shutdown(ctx)
}

// Db returns the database instance.
func (w *Way) Db() *DB {
	w.Logger.Println("Database instance returned")
	return w.db
}

// GetEnv retrieves the environment variable value or a default value if not set
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Environment variables
const (
	// Environment variables for database connection
	envDBDriver   = "WAY_DB_DRIVER"
	envDBUsePGX   = "WAY_DB_USE_PGX"
	envDBDSN      = "WAY_DB_DSN"
	envDBUser     = "WAY_DB_USER"
	envDBPassword = "WAY_DB_PASSWORD"
	envDBHost     = "WAY_DB_HOST"
	envDBPort     = "WAY_DB_PORT"
	envDBName     = "WAY_DB_NAME"
	// Environment variables for session management
	envStoreName          = "WAY_DEFAULT_STORE_NAME"
	envStoreEncryptionKey = "WAY_DEFAULT_STORE_ENCRYPTION_KEY"
	// Environment variables for cookie management
	envCookieName              = "WAY_DEFAULT_COOKIE_NAME"
	envCookieEncryptionKey     = "WAY_DEFAULT_COOKIE_ENCRYPTION_KEY"
	envCookieAuthenticationKey = "WAY_DEFAULT_COOKIE_AUTHENTICATION_KEY"
	// Environment variables for DefaultLogger
	envDefaultLogger = "WAY_DEFAULT_LOGGER"
)

// defaultLogger returns a new logger that writes to os.Stdout.
func defaultLogger() *log.Logger {
	return log.New(os.Stdout, GetEnv(envDefaultLogger, "WAY_INFO")+": ", log.LstdFlags)
}

func loggingMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Printf("method=%s path=%s duration=%s", r.Method, r.URL.Path, time.Since(start))
	})
}

// useDefaultSession checks if the default session should be used.
func useDefaultSession() bool {
	return GetEnv(envStoreEncryptionKey, "") != "" ||
		(GetEnv(envCookieEncryptionKey, "") != "" && GetEnv(envCookieAuthenticationKey, "") != "")
}

// setSessionDefaults sets the default values for a Session object.
func setSessionDefaults(s *Session) {
	s.defaultStore = getDefaultStoreName()
	s.defaultCookie = getDefaultCookieName()
	if key := GetEnv(envStoreEncryptionKey, ""); key != "" {
		s.stores[s.defaultStore] = sessions.NewCookieStore([]byte(key))
	}
	if encKey := GetEnv(envCookieEncryptionKey, ""); encKey != "" {
		authKey := GetEnv(envCookieAuthenticationKey, "")
		s.cookies[s.defaultCookie] = securecookie.New([]byte(encKey), []byte(authKey))
	}
}

// getDefaultCookieName returns the default cookie name.
func getDefaultCookieName() string {
	return GetEnv(envCookieName, "way")
}

// getDefaultStoreName returns the default store name.
func getDefaultStoreName() string {
	return GetEnv(envStoreName, "way")
}

// getEncryptionKey retrieves the encryption key used for cookie encryption.
func getEncryptionKey() []byte {
	key := GetEnv(envCookieEncryptionKey, "")
	if key == "" {
		log.Fatalf("%s is required", envCookieEncryptionKey)
	}
	return []byte(key)
}

// getAuthenticationKey retrieves the authentication key from the environment variable.
func getAuthenticationKey() []byte {
	key := GetEnv(envCookieAuthenticationKey, "")
	if key == "" {
		log.Fatalf("%s is required", envCookieAuthenticationKey)
	}
	return []byte(key)
}

// getStoreEncryptionKey retrieves the encryption key used for storing data in the application's store.
func getStoreEncryptionKey() []byte {
	key := GetEnv(envStoreEncryptionKey, "")
	if key == "" {
		log.Fatalf("%s is required", envStoreEncryptionKey)
	}
	return []byte(key)
}
