package auth

import (
	"context"
	"log"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const userContextKey = contextKey("user")

// UserClaims holds the JWT claims
type UserClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

// AddUserToContext adds the user claims from the JWT to the context
func AddUserToContext(ctx context.Context, claims *jwt.StandardClaims) context.Context {
	// Convert the UserID from string to uint
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		log.Printf("Failed to parse user ID from token claims: %v", err)
		return ctx // Return the original context if parsing fails
	}

	userClaims := UserClaims{
		UserID: uint(userID), // Cast uint64 to uint
	}

	log.Printf("UserID added to context: %d", userClaims.UserID) // Log for debugging
	return context.WithValue(ctx, userContextKey, userClaims)
}

// GetUserFromContext retrieves the user claims from the context
func GetUserFromContext(ctx context.Context) (*UserClaims, bool) {
	userClaims, ok := ctx.Value(userContextKey).(UserClaims)
	return &userClaims, ok
}
