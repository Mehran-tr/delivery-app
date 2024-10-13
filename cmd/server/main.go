package main

import (
	"go-delivery-app/internal/db"
	"go-delivery-app/internal/notifications"
	"go-delivery-app/internal/routes"
	"log"
	"net/http"
)

func main() {
	// Connect to the database
	if err := db.ConnectDatabase(); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Run migrations
	db.RunMigrations(db.DB)

	// Initialize routes from the routes package
	router := routes.InitializeRoutes()

	// Start RabbitMQ consumers to process notifications
	go notifications.ConsumeNotifications("notifications_sender_queue")
	go notifications.ConsumeNotifications("notifications_motorbike_queue")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}
