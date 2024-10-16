package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-delivery-app/internal/auth"
	"go-delivery-app/internal/db"
	"go-delivery-app/internal/models"
	"go-delivery-app/internal/notifications"
	"net/http"
	"time"
)

// CreateParcel allows a sender to create a new parcel, automatically setting the SenderID from the authenticated user
func CreateParcel(w http.ResponseWriter, r *http.Request) {
	var parcel models.Parcel

	// Decode the JSON request body into a Parcel struct
	err := json.NewDecoder(r.Body).Decode(&parcel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the authenticated user's ID (SenderID) from the request context
	userClaims, ok := auth.GetUserFromContext(r.Context())
	if !ok || userClaims.UserID == 0 {
		http.Error(w, "Could not determine authenticated user", http.StatusUnauthorized)
		return
	}

	// Ensure that required fields are provided
	if parcel.PickupAddress == "" || parcel.DropoffAddress == "" || parcel.Latitude == 0 || parcel.Longitude == 0 {
		http.Error(w, "PickupAddress, DropoffAddress, Latitude, and Longitude are required", http.StatusBadRequest)
		return
	}

	// Set the SenderID automatically from the authenticated user
	parcel.SenderID = userClaims.UserID
	parcel.Status = "Created" // Set initial status as "Created"

	// Save the parcel to the database
	result := db.DB.Create(&parcel)
	if result.Error != nil {
		http.Error(w, "Failed to save parcel", http.StatusInternalServerError)
		return
	}
	// Send a notification to the sender using RabbitMQ
	message := "Your parcel has been created successfully!"
	notifications.PublishNotification("notifications_sender_queue", parcel.SenderID, message)

	// Set response header to application/json and return the created parcel
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(parcel)
}

// GetParcelStatus allows a sender to check the status of a parcel, but only if it's their own
func GetParcelStatus(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user (sender) from the request context
	userClaims, ok := auth.GetUserFromContext(r.Context())
	if !ok || userClaims.UserID == 0 {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	// Retrieve the parcel ID from the request path
	parcelID := mux.Vars(r)["id"]
	var parcel models.Parcel

	// Retrieve the parcel from the database by ID
	db.DB.First(&parcel, parcelID)

	// If the parcel doesn't exist, return a 404 error
	if parcel.ID == 0 {
		http.Error(w, "Parcel not found", http.StatusNotFound)
		return
	}

	// Check if the parcel belongs to the authenticated sender
	if parcel.SenderID != userClaims.UserID {
		http.Error(w, "You are not authorized to view this parcel", http.StatusForbidden)
		return
	}

	// Set the response header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Return the parcel status as a single JSON object
	json.NewEncoder(w).Encode(parcel)
}

// CancelParcel allows senders or motorbikes to cancel a parcel
func CancelParcel(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user's claims (either sender or motorbike)
	userClaims, ok := auth.GetUserFromContext(r.Context())
	if !ok || userClaims.UserID == 0 {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	// Retrieve the parcel ID from the URL
	parcelID := mux.Vars(r)["id"]
	var parcel models.Parcel

	// Retrieve the parcel from the database by ID
	db.DB.First(&parcel, parcelID)

	// If the parcel doesn't exist, return a 404 error
	if parcel.ID == 0 {
		http.Error(w, "Parcel not found", http.StatusNotFound)
		return
	}

	// Check if the parcel is already delivered, which cannot be canceled
	if parcel.Status == "Delivered" {
		http.Error(w, "Delivered parcels cannot be canceled", http.StatusBadRequest)
		return
	}

	// Check if the user is authorized to cancel the parcel
	// Sender can only cancel their own parcels and if the status is "Created"
	if userClaims.Role == "sender" && parcel.SenderID != userClaims.UserID {
		http.Error(w, "You can only cancel your own parcels", http.StatusForbidden)
		return
	}

	// Motorbike can cancel parcels they have picked up, but only if the status is "Picked up"
	if userClaims.Role == "motorbike" && (parcel.MotorbikeID == nil || *parcel.MotorbikeID != userClaims.UserID) {
		http.Error(w, "You can only cancel parcels you have picked up", http.StatusForbidden)
		return
	}

	// Update the parcel's status to "Canceled"
	cancelTime := time.Now()
	parcel.Status = "Canceled"
	parcel.CanceledAt = &cancelTime

	// Save the updated parcel back to the database
	result := db.DB.Save(&parcel)
	if result.Error != nil {
		http.Error(w, "Failed to cancel the parcel", http.StatusInternalServerError)
		return
	}

	// Send notifications to both the sender and the motorbike (if applicable)
	if parcel.SenderID != 0 {
		notifications.PublishNotification("notifications_sender_queue", parcel.SenderID, "Your parcel has been canceled")
	}
	if parcel.MotorbikeID != nil {
		notifications.PublishNotification("notifications_motorbike_queue", *parcel.MotorbikeID, "The parcel you picked up has been canceled")
	}

	// Send the updated parcel as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Parcel has been canceled"})
}
