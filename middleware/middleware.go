package middleware

import (
	"net/http"
)

type Middleware interface {
	// Middleware is a function that wraps an http.Handler.
	Middleware(next http.Handler) http.Handler
}
