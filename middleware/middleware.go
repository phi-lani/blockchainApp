package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/phi-lani/blockchainApp/utils"
)

// LoggingMiddleware - Logs each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// AuthenticationMiddleware - Validates JWT for protected routes
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token from header
		token := strings.Split(authHeader, "Bearer ")[1]

		// Validate token
		if !utils.ValidateJWT(token) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
