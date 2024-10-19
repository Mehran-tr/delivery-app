package models

import "time"

// User represents the structure of users (senders, motorbikes, and admins).
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"`
}

type Parcel struct {
	ID                   uint       `gorm:"primaryKey"`
	SenderID             uint       `json:"SenderID"`
	PickupAddress        string     `json:"PickupAddress"`
	DropoffAddress       string     `json:"DropoffAddress"`
	Latitude             float64    `json:"Latitude"`
	Longitude            float64    `json:"Longitude"`
	Status               string     `json:"Status"`
	PickupTime           *time.Time `json:"PickupTime"`
	DeliveryTime         *time.Time `json:"DeliveryTime"`
	MotorbikeID          *uint      `json:"MotorbikeID"`
	SenderDescription    *string    `json:"SenderDescription"`    // Nullable field
	MotorbikeDescription *string    `json:"MotorbikeDescription"` // Nullable field
	CanceledAt           *time.Time `json:"canceled_at"`          // Nullable field
}

// Notification represents a notification to be sent to a user
type Notification struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `json:"UserID"`    // The recipient of the notification
	Message   string    `json:"Message"`   // The notification message
	CreatedAt time.Time `json:"CreatedAt"` // Timestamp for the notification
	Read      bool      `json:"Read"`      // Whether the notification has been read
}

// Rating represents the rating given by a sender to a motorbike after parcel delivery
type Rating struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	SenderID    uint      `json:"sender_id"`    // ID of the sender who gave the rating
	MotorbikeID uint      `json:"motorbike_id"` // ID of the motorbike being rated
	ParcelID    uint      `json:"parcel_id"`    // ID of the parcel related to the rating
	Rating      int       `json:"rating"`       // Rating value (1 to 5)
	CreatedAt   time.Time `json:"created_at"`   // Timestamp when the rating was created
}
