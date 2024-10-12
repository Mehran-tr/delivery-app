package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging logs each request made to the server
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log the details of the request
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the completion time and status
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}
