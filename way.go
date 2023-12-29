package way

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Way struct {
	startupMutex sync.RWMutex
	router       *mux.Router
	sql          *sql.DB
	Server       *http.Server
	Listener     net.Listener
}

type HandlerFunc func(*Context)
type MiddlewareFunc func(HandlerFunc) HandlerFunc

func New() *Way {
	return &Way{
		router: mux.NewRouter(),
		Server: new(http.Server),
	}
}

// adaptHandler adapts a `HandlerFunc` to `http.HandlerFunc`.
func adaptHandler(db *sql.DB, handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(db, w, r)
		handler(ctx)
	}
}

// handleFuncWithMethod registers a new route with a matcher for the URL path and the HTTP method.
func (w *Way) handleFuncWithMethod(path string, handler HandlerFunc, method string) {
	w.router.HandleFunc(path, adaptHandler(w.sql, handler)).Methods(method)
}

// Use adds a middleware to the middleware stack.
func (w *Way) Use(middleware ...MiddlewareFunc) {
	w.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w.sql, wr, r)

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

// HandleFunc registers a new route with a matcher for the URL path.
func (w *Way) HandleFunc(path string, handler HandlerFunc) {
	w.router.HandleFunc(path, adaptHandler(w, handler))
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

// newListener creates a new net.Listener.
func newListener(network, address string) (net.Listener, error) {
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	return l, nil
}
