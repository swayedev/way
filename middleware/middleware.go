package middleware

import (
	"log"
	"net/http"
)

type Middleware interface {
	// Middleware is a function that wraps an http.Handler.
	Middleware(logger *log.Logger, next http.Handler) http.Handler
}
