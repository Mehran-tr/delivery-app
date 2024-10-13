package handlers

import (
	"go-delivery-app/internal/auth"
	"go-delivery-app/internal/db"
	"go-delivery-app/internal/models"
	"net/http"

	"github.com/gorilla/mux"
)

// MarkNotificationAsRead allows a user to mark a notification as read
func MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user's claims
	userClaims, ok := auth.GetUserFromContext(r.Context())
	if !ok || userClaims.UserID == 0 {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	// Get notification ID from the URL
	notificationID := mux.Vars(r)["id"]
	var notification models.Notification

	// Find the notification
	db.DB.First(&notification, notificationID)

	// Check if the notification exists and belongs to the user
	if notification.ID == 0 || notification.UserID != userClaims.UserID {
		http.Error(w, "Notification not found or unauthorized", http.StatusForbidden)
		return
	}

	// Mark the notification as read
	notification.Read = true
	db.DB.Save(&notification)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification marked as read"))
}
