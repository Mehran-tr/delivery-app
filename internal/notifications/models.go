package notifications

// NotificationMessage defines the structure of the notification message sent via RabbitMQ
type NotificationMessage struct {
	UserID  uint   `json:"user_id"`
	Message string `json:"message"`
}
