# API Freeze for Way v1.0.0-rc1

This document outlines the public API surface and stability guarantees for Way v1.0.0-rc1. 

## Frozen APIs

The following are considered stable and changes will follow semantic versioning rules:

### Way Core Package (`github.com/swayedev/way`)

**Type: Way**
- `func New() *Way` – creates a new Way instance with safe server defaults
- `func (w *Way) SetLogger(logger *log.Logger)` – replaces the logger
- `func (w *Way) SetRouter(router *mux.Router)` – replaces the router
- `func (w *Way) SetServer(server *http.Server)` – replaces the HTTP server
- `func (w *Way) SetListener(listener net.Listener)` – replaces the network listener
- `func (w *Way) SetHTTPClient(client *http.Client)` – sets the outbound client used by context proxy helpers
- `func (w *Way) SetDB(db *DB)` – sets the database connection
- `func (w *Way) SetSession(s *Session)` – sets the session manager
- `func (w *Way) Use(middleware ...MiddlewareFunc)` – registers middleware
- `func (w *Way) HandleFunc(path string, handler HandlerFunc)` – registers a route
- `func (w *Way) GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD(path string, handler HandlerFunc)` – HTTP method helpers
- `func (w *Way) Start(address string) error` – starts the server
- `func (w *Way) Close() error` – immediately stops the server
- `func (w *Way) Shutdown(ctx context.Context) error` – gracefully shuts down the server
- `func (w *Way) Db() *DB` – retrieves the database connection
- `func (w *Way) Log() *log.Logger` – retrieves the logger
- `func (w *Way) InitDBFromConfig() error` – initializes database from environment variables

**Type: Context**
- JSON/XML/String/HTML/Data/Image/Redirect/Header/Cookie/Status response helpers
- `func (c *Context) ProxyMedia(url string)` – proxy media to the client using the configured HTTP client
- `func (c *Context) GetSession(name string) sessions.Store` – compatibility session-store lookup
- `func (c *Context) GetSessionE(name string) (sessions.Store, error)` – error-returning session-store lookup

**Types: HandlerFunc, MiddlewareFunc**
- Route handlers and middleware chainable wrappers

**Type: Session**
- Session and secure cookie store management with named stores and cookies
- `func (s *Session) Store(name string) sessions.Store` – retrieve a session store
- `func (s *Session) StoreE(name string) (sessions.Store, error)` – retrieve a session store with an explicit error
- `func (s *Session) Cookie(name string) *securecookie.SecureCookie` – retrieve a secure cookie
- `func (s *Session) CookieE(name string) (*securecookie.SecureCookie, error)` – retrieve a secure cookie with an explicit error

**Type: DB**
- `func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error)`
- `func (d *DB) QueryRow(query string, args ...interface{}) *sql.Row`
- `func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error)`
- SQL and pgx connection management via `SetDB` and `SetDBConnection`

### Way Crypto Package (`github.com/swayedev/way/crypto`)

- `func Encrypt(data []byte, passphrase string) (string, error)` – encrypt with passphrase, return hex string
- `func Decrypt(encrypted string, passphrase string) ([]byte, error)` – decrypt hex string with passphrase
- `func HashString(data string) [32]byte` – SHA3-256 hash
- `func HashByte(data []byte) [32]byte` – SHA3-256 hash
- `func HashStringToString(data string) string` – SHA3-256 hash as hex
- `func GenerateRandomKey(length int) ([]byte, error)` – generate random bytes with error return

### Way Database Package (`github.com/swayedev/way/database`)

- `func CheckDriver(driver string) string` – validate driver name
- `func CheckDSN(driver, dsn, name, host, port, user, password string) string` – validate DSN
- `func DriverImportHint(driver string) string` – returns optional driver adapter guidance
- `func SQLConnect(driver, dsn string) (*sql.DB, error)` – open SQL connection
- `func PGXConnect(dsn string) (*pgx.Conn, error)` – open PGX connection (PostgreSQL)

### Way Database Driver Adapter Packages

Applications should import the adapter package for every `database/sql` driver they use:

- `github.com/swayedev/way/database/drivers/mysql`
- `github.com/swayedev/way/database/drivers/pgx`
- `github.com/swayedev/way/database/drivers/sqlite`
- `github.com/swayedev/way/database/drivers/sqlserver`
- `github.com/swayedev/way/database/drivers/godror`

## Breaking Changes from Pre-1.0

The following are intentional breaking changes for v1.0.0-rc1:

1. **Error returns**: `getEncryptionKey()`, `getAuthenticationKey()`, `getStoreEncryptionKey()`, and `GenerateRandomKey()` now return `error` as a second return value instead of calling `log.Fatalf` or silently returning `nil`.
2. **Server defaults**: `New()` now sets safe `http.Server` timeouts by default: `ReadHeaderTimeout: 5s`, `ReadTimeout: 15s`, `WriteTimeout: 15s`, `IdleTimeout: 30s`.
3. **ASCII art**: Server startup ASCII art is now gated behind the `WAY_LOG_ASCII_ART=true` environment variable and off by default.
4. **Database drivers**: Core Way no longer blank-imports SQL drivers. Import the adapter package for your database driver.
5. **Session errors**: Missing session stores and secure-cookie configurations now return errors from `*E` methods and from encrypted-cookie helpers instead of panicking.

## Deprecations (Compatibility Kept)

The non-error-returning `Session.Store`, `Session.Cookie`, `Session.DefaultSession`, `Session.DefaultCookie`, and `Context.GetSession` methods remain for compatibility. New code should prefer the `*E` variants where available.

## Future Stability

Changes beyond v1.0.0-rc1 will follow Semantic Versioning:
- **Patch** (v1.0.1): bug fixes, security patches, internal optimizations
- **Minor** (v1.1.0): new methods, new drivers, backward-compatible enhancements
- **Major** (v2.0.0): breaking API changes only

## Migration from Pre-rc1

If you are upgrading from a pre-rc1 version:

1. Update error handling for the four functions that now return errors (see Breaking Changes above).
2. Ensure `WAY_LOG_ASCII_ART=true` is set if you relied on startup ASCII art output.
3. Review `way.New()` and custom `http.Server` configuration if you were managing timeouts manually.
4. Import a database driver adapter package before calling SQL connection helpers.

Refer to [MIGRATION.md](MIGRATION.md) for specific upgrade paths.
