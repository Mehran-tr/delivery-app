package services

import (
	"context"
	"errors"
	"go-delivery-app/internal/auth"
)

// AuthenticatedUser represents the details of the authenticated user.
type AuthenticatedUser struct {
	UserID uint
	Role   string
}

// GetAuthenticatedUser retrieves the authenticated user's details from the request context.
func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, error) {
	// Retrieve the user claims from the context
	userClaims, ok := auth.GetUserFromContext(ctx)
	if !ok || userClaims.UserID == 0 {
		return nil, errors.New("could not determine authenticated user")
	}

	// Return the authenticated user details
	return &AuthenticatedUser{
		UserID: userClaims.UserID,
		Role:   userClaims.Role,
	}, nil
}
