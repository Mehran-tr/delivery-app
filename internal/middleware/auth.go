package middleware

import (
	"go-delivery-app/internal/auth"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// JWTMiddleware checks for a valid JWT token in the Authorization header
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		// Extract the token from the "Bearer <token>" format
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the JWT token
		jwtSecret := os.Getenv("JWT_SECRET") // Ensure the same JWT_SECRET is used
		token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil // Use the same secret used for signing
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims and add them to the context
		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			r = r.WithContext(auth.AddUserToContext(r.Context(), claims))
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
	})
}
