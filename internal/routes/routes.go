package routes

import (
	"go-delivery-app/internal/handlers"
	"go-delivery-app/internal/middleware"

	"github.com/gorilla/mux"
)

// InitializeRoutes sets up all the routes for the application
func InitializeRoutes() *mux.Router {
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	// Protected routes with JWT middleware
	senderRoutes := router.PathPrefix("/sender").Subrouter()
	senderRoutes.Use(middleware.JWTMiddleware)
	senderRoutes.HandleFunc("/parcel", handlers.CreateParcel).Methods("POST")
	senderRoutes.HandleFunc("/parcel/{id}", handlers.GetParcelStatus).Methods("GET")

	motorbikeRoutes := router.PathPrefix("/motorbike").Subrouter()
	motorbikeRoutes.Use(middleware.JWTMiddleware)
	motorbikeRoutes.HandleFunc("/parcels", handlers.ListParcels).Methods("GET")
	motorbikeRoutes.HandleFunc("/parcel/{id}/pickup", handlers.PickParcel).Methods("POST")
	motorbikeRoutes.HandleFunc("/parcel/{id}/update", handlers.UpdateParcelStatus).Methods("PUT")

	adminRoutes := router.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middleware.JWTMiddleware)
	adminRoutes.HandleFunc("/parcels", handlers.GetAllParcels).Methods("GET")
	adminRoutes.HandleFunc("/users", handlers.GetUsers).Methods("GET")

	return router
}
