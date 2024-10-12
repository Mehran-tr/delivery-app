package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateToken creates a JWT token for a user
func GenerateToken(userID uint, role string) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(72 * time.Hour).Unix(), // Token expiration
		Subject:   strconv.Itoa(int(userID)),             // Store the userID as string in the token
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	jwtSecret := os.Getenv("JWT_SECRET") // Make sure JWT_SECRET is set correctly
	return token.SignedString([]byte(jwtSecret))
}
