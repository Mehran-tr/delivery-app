package notifications

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// PublishNotification sends a notification message to the RabbitMQ queue
func PublishNotification(queueName string, userID uint, message string) error {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		queueName, // name of the queue
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	// Create the notification message with userID and message
	notification := NotificationMessage{
		UserID:  userID,
		Message: message,
	}

	// Convert the notification to JSON
	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %v", err)
	}

	// Publish the message to the queue
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	log.Printf(" [x] Sent %s", message)
	return nil
}
