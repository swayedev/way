package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// RecoverMiddleware recovers from panics and logs the error.
func RecoverMiddleware(logger *log.Logger, next http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Printf("panic: %v\n%s", err, debug.Stack())
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
