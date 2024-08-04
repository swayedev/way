package way

import (
	"context"
	"encoding/base64"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/swayedev/way/database"
	"golang.org/x/exp/rand"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// Config holds all configuration parameters.
type Config struct {
	DBDriver                string
	DBUsePGX                bool
	DBDSN                   string
	DBUser                  string
	DBPassword              string
	DBHost                  string
	DBPort                  string
	DBName                  string
	StoreName               string
	StoreEncryptionKey      string
	CookieName              string
	CookieEncryptionKey     string
	CookieAuthenticationKey string
	DefaultLogger           string
	Crypto                  *Crypto
}

// Way is the main framework structure.
type Way struct {
	startupMutex sync.RWMutex
	db           *DB
	router       *mux.Router
	sessions     *Session
	Server       *http.Server
	Listener     net.Listener
	Logger       *log.Logger
	Config       *Config
	Crypto       *Crypto
}

// HandlerFunc is a function type that represents a handler for a request.
type HandlerFunc func(*Context)

// MiddlewareFunc represents a function that takes a HandlerFunc and returns a modified HandlerFunc.
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// New creates a new instance of Way.
func New(config *Config) *Way {
	genKeys := false
	if config == nil {
		config = defaultConfig()
		genKeys = true
	}
	logger := defaultLogger(config.DefaultLogger)
	if genKeys {
		config.StoreEncryptionKey = generateRandomKey(32)
		config.CookieEncryptionKey = generateRandomKey(32)
		config.CookieAuthenticationKey = generateRandomKey(32)
		logger.Println("Generated new encryption keys for store and cookies")
	}
	sessions := NewSession(logger)
	if useDefaultSession(config) {
		setSessionDefaults(sessions, config)
	}
	return &Way{
		router:   mux.NewRouter(),
		Server:   &http.Server{},
		sessions: sessions,
		Logger:   logger,
		Config:   config,
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

// SetCrypto sets the Crypto object.
func (w *Way) SetCrypto(c *Crypto) {
	w.Crypto = c
}

// Log returns the logger.
func (w *Way) Log() *log.Logger {
	if w.Logger != nil {
		return w.Logger
	}
	return defaultLogger(w.Config.DefaultLogger)
}

// SetDB sets the database connection for the Way object.
func (w *Way) SetDB(db *DB) {
	w.db = db
}

func (w *Way) initDBConnection(driver, dsn string) error {
	if w.Config.DBUsePGX {
		db, _, err := database.PGXConnect(database.PGXConfig{DSN: dsn, UsePooling: false})
		if err != nil {
			return err
		}
		w.db = &DB{pgx: db, UsePgx: true, Driver: "pgx"}
		return nil
	}
	db, err := database.SQLConnect(database.SQLConfig{Driver: driver, DSN: dsn, UsePooling: false})
	if err != nil {
		return err
	}
	w.db = &DB{sql: db, UsePgx: false, Driver: driver}
	return nil
}

// InitDBFromConfig initializes the database connection from the Config struct.
func (w *Way) InitDBFromConfig() error {
	driver, err := database.CheckDriver(w.Config.DBDriver)
	if err != nil {
		return err
	}
	dsn, err := database.CheckDSN(database.DriverConfig{
		Driver: w.Config.DBDriver,
		DSN:    w.Config.DBDSN,
		DBName: w.Config.DBName,
		DBHost: w.Config.DBHost,
		DBPort: w.Config.DBPort,
		DBUser: w.Config.DBUser,
		DBPass: w.Config.DBPassword,
	})
	if err != nil {
		return err
	}
	return w.initDBConnection(driver, dsn)
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

// Use adds middleware to the middleware stack.
func (w *Way) Use(middleware ...MiddlewareFunc) {
	for _, m := range middleware {
		w.router.Use(adaptMiddleware(w, m))
	}
}

// adaptMiddleware adapts a MiddlewareFunc to mux.MiddlewareFunc.
func adaptMiddleware(w *Way, m MiddlewareFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
			ctx := NewContext(wr, r, w.db, w.sessions, w.Logger, w.Crypto)
			middlewareHandler := m(func(c *Context) {
				next.ServeHTTP(wr, r)
			})
			middlewareHandler(ctx)
		})
	}
}

// adaptHandler adapts a HandlerFunc to http.HandlerFunc.
func adaptHandler(wa *Way, handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r, wa.db, wa.sessions, wa.Logger, wa.Crypto)
		handler(ctx)
	}
}

// handleFuncWithMethod registers a new route with a matcher for the URL path and the HTTP method.
func (w *Way) handleFuncWithMethod(path string, handler HandlerFunc, method string) {
	w.Logger.Printf("Registering route %s", path)
	w.router.HandleFunc(path, adaptHandler(w, handler)).Methods(method)
}

// HandleFunc registers a new route with a matcher for the URL path.
func (w *Way) HandleFunc(path string, handler HandlerFunc) {
	w.Logger.Printf("Registering route %s", path)
	w.router.HandleFunc(path, adaptHandler(w, handler))
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
	w.Logger.Println(`
	___________________________________
	       __        ______   __
	       \ \      / /  \ \ / /
	        \ \ /\ / / /\ \ V / 
	         \ V  V / /__\ | |  
	          \_/\_/_/----\|_|  
	------------------------------------
		Server is running an port: %s
	------------------------------------
	`, address)
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

// GetEnv retrieves the environment variable value or a default value if not set.
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
func defaultLogger(defaultLogger string) *log.Logger {
	return log.New(os.Stdout, defaultLogger+": ", log.LstdFlags)
}

func loggingMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Printf("method=%s path=%s duration=%s", r.Method, r.URL.Path, time.Since(start))
	})
}

// useDefaultSession checks if the default session should be used.
func useDefaultSession(config *Config) bool {
	return config.StoreEncryptionKey != "" ||
		(config.CookieEncryptionKey != "" && config.CookieAuthenticationKey != "")
}

// setSessionDefaults sets the default values for a Session object.
func setSessionDefaults(s *Session, config *Config) {
	s.defaultStore = config.StoreName
	s.defaultCookie = config.CookieName
	if key := config.StoreEncryptionKey; key != "" {
		s.stores[s.defaultStore] = sessions.NewCookieStore([]byte(key))
	}
	if encKey := config.CookieEncryptionKey; encKey != "" {
		authKey := config.CookieAuthenticationKey
		s.cookies[s.defaultCookie] = securecookie.New([]byte(encKey), []byte(authKey))
	}
}

// defaultConfig returns a Config struct with default values.
func defaultConfig() *Config {
	return &Config{
		DBDriver:                GetEnv(envDBDriver, "postgres"),
		DBUsePGX:                GetEnv(envDBUsePGX, "true") == "true",
		DBDSN:                   GetEnv(envDBDSN, ""),
		DBUser:                  GetEnv(envDBUser, "user"),
		DBPassword:              GetEnv(envDBPassword, "password"),
		DBHost:                  GetEnv(envDBHost, "localhost"),
		DBPort:                  GetEnv(envDBPort, "5432"),
		DBName:                  GetEnv(envDBName, "dbname"),
		StoreName:               GetEnv(envStoreName, "default"),
		StoreEncryptionKey:      GetEnv(envStoreEncryptionKey, ""),
		CookieName:              GetEnv(envCookieName, "way"),
		CookieEncryptionKey:     GetEnv(envCookieEncryptionKey, ""),
		CookieAuthenticationKey: GetEnv(envCookieAuthenticationKey, ""),
		DefaultLogger:           GetEnv(envDefaultLogger, "WAY_INFO"),
	}
}

// generateRandomKey generates a random key of the given length.
func generateRandomKey(length int) string {
	key := make([]byte, length)
	if _, err := rand.Read(key); err != nil {
		log.Fatalf("Failed to generate random key: %v", err)
	}
	return base64.StdEncoding.EncodeToString(key)
}
