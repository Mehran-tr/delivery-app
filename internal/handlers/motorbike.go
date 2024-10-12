package handlers

import (
	"encoding/json"
	"go-delivery-app/internal/auth"
	"go-delivery-app/internal/db"
	"go-delivery-app/internal/models"
	"go-delivery-app/internal/services"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// ListParcels allows motorbikes to see available parcels (returns individual JSON objects instead of an array)
func ListParcels(w http.ResponseWriter, r *http.Request) {
	var parcels []models.Parcel
	db.DB.Where("motorbike_id IS NULL").Find(&parcels)

	// Set the response content type
	w.Header().Set("Content-Type", "application/json")

	// Iterate over parcels and write each one as an individual JSON object
	for i, parcel := range parcels {
		if i > 0 {
			// Write a comma between objects if it's not the first parcel
			w.Write([]byte(","))
		}
		json.NewEncoder(w).Encode(parcel)
	}
}

// Request body struct to capture motorbike description
type PickParcelRequest struct {
	MotorbikeDescription string `json:"MotorbikeDescription"` // PascalCase for JSON field
}

// PickParcel allows motorbikes to pick up a parcel by its ID and add motorbike description
func PickParcel(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user's claims (motorbike)
	userClaims, ok := auth.GetUserFromContext(r.Context())
	if !ok || userClaims.UserID == 0 {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	// Check if the authenticated user is a motorbike
	if userClaims.Role == "sender" {
		http.Error(w, "Only motorbikes can pick parcels", http.StatusForbidden)
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

	// If the parcel is already picked up or delivered, it cannot be picked again
	if parcel.Status == "Picked up" || parcel.Status == "Delivered" {
		http.Error(w, "Parcel already picked up or delivered", http.StatusBadRequest)
		return
	}

	// Parse the request body to get motorbike description
	var reqBody PickParcelRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the parcel's status, motorbike ID, and motorbike description
	pickupTime := time.Now()
	parcel.Status = "Picked up"
	parcel.PickupTime = &pickupTime
	parcel.MotorbikeID = &userClaims.UserID                     // Assign the motorbike ID from the authenticated user
	parcel.MotorbikeDescription = &reqBody.MotorbikeDescription // Assign motorbike description (pointer)

	// Save the updated parcel back to the database
	result := db.DB.Save(&parcel)
	if result.Error != nil {
		http.Error(w, "Failed to update parcel", http.StatusInternalServerError)
		return
	}

	// Send the updated parcel as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parcel)
}

// UpdateParcelStatus allows motorbikes to update the status of a parcel to "Delivered"
func UpdateParcelStatus(w http.ResponseWriter, r *http.Request) {
	// Use the service to get the authenticated user's details
	user, err := services.GetAuthenticatedUser(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Check if the authenticated user is a motorbike
	if user.Role == "sender" {
		http.Error(w, "Only motorbikes can update parcel status", http.StatusForbidden)
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

	// Check if the parcel has been picked up before updating to "Delivered"
	if parcel.Status != "Picked up" {
		http.Error(w, "Parcel has not been picked up yet", http.StatusBadRequest)
		return
	}

	// Update the parcel's status to "Delivered"
	deliveryTime := time.Now()
	parcel.Status = "Delivered"
	parcel.DeliveryTime = &deliveryTime

	// Save the updated parcel status to the database
	result := db.DB.Save(&parcel)
	if result.Error != nil {
		http.Error(w, "Failed to update parcel status", http.StatusInternalServerError)
		return
	}

	// Send the updated parcel status as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Parcel marked as delivered"})
}
