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
	senderRoutes.HandleFunc("/parcel/{id}/cancel", handlers.CancelParcel).Methods("POST")

	motorbikeRoutes := router.PathPrefix("/motorbike").Subrouter()
	motorbikeRoutes.Use(middleware.JWTMiddleware)
	motorbikeRoutes.HandleFunc("/parcels", handlers.ListParcels).Methods("GET")
	motorbikeRoutes.HandleFunc("/parcel/{id}/pickup", handlers.PickParcel).Methods("POST")
	motorbikeRoutes.HandleFunc("/parcel/{id}/update", handlers.UpdateParcelStatus).Methods("PUT")
	motorbikeRoutes.HandleFunc("/parcel/{id}/cancel", handlers.CancelParcel).Methods("POST")

	adminRoutes := router.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middleware.JWTMiddleware)
	adminRoutes.HandleFunc("/parcels", handlers.GetAllParcels).Methods("GET")
	adminRoutes.HandleFunc("/users", handlers.GetUsers).Methods("GET")

	// Notification routes
	notificationRoutes := router.PathPrefix("/notifications").Subrouter()
	notificationRoutes.Use(middleware.JWTMiddleware)
	notificationRoutes.HandleFunc("", handlers.GetNotifications).Methods("GET")
	notificationRoutes.HandleFunc("/{id}/read", handlers.MarkNotificationAsRead).Methods("PUT")

	return router
}
