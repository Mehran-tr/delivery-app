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
}
