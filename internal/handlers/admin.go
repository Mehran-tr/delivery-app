package handlers

import (
	"encoding/json"
	"go-delivery-app/internal/db"
	"go-delivery-app/internal/models"
	"net/http"
)

// GetAllParcels allows admin to see all parcels
func GetAllParcels(w http.ResponseWriter, r *http.Request) {
	var parcels []models.Parcel
	db.DB.Find(&parcels)

	// Start the response header with application/json content-type
	w.Header().Set("Content-Type", "application/json")

	// Iterate over users and write each one as an individual JSON object
	for i, parcel := range parcels {
		if i > 0 {
			// Write a comma between objects if it's not the first parcel
			w.Write([]byte(","))
		}

		// Encode and write the current parcel object
		json.NewEncoder(w).Encode(parcel)
	}

}

// GetUsers allows admin to see all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	db.DB.Find(&users)

	// Start the response header with application/json content-type
	w.Header().Set("Content-Type", "application/json")

	// Iterate over users and write each one as an individual JSON object
	for i, user := range users {
		if i > 0 {
			// Write a comma between objects if it's not the first user
			w.Write([]byte(","))
		}

		// Encode and write the current user object
		json.NewEncoder(w).Encode(user)
	}
}
