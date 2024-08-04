package middleware

// LoggerMiddleware logs each request with its method, path, and duration.
// func LoggerMiddleware(next http.Handler) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			start := time.Now()
// 			next.ServeHTTP(w, r)
// 			logger.Printf("method=%s path=%s duration=%s", r.Method, r.URL.Path, time.Since(start))
// 		})
// 	}
// }
