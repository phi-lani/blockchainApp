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
		// Get the Authorization header
		tokenString := r.Header.Get("Authorization")

		// Check if the Authorization header is provided
		if tokenString == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Extract the token from the "Bearer <token>" string
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Validate the token
		token, err := utils.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
