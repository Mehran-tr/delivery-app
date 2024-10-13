package handlers

import (
	"encoding/json"
	"go-delivery-app/internal/auth"
	"go-delivery-app/internal/db"
	"go-delivery-app/internal/models"
	"net/http"
)

// GetNotifications allows a user to retrieve their notifications
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user's claims
	userClaims, ok := auth.GetUserFromContext(r.Context())
	if !ok || userClaims.UserID == 0 {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	var notifications []models.Notification
	db.DB.Where("user_id = ?", userClaims.UserID).Find(&notifications)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}
