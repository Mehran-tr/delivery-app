package notifications

import (
	"encoding/json"
	"go-delivery-app/internal/db"
	"go-delivery-app/internal/models"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// ConsumeNotifications listens to the RabbitMQ queue and processes notifications
func ConsumeNotifications(queueName string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name of the queue
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer tag
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			// Unmarshal the JSON message into NotificationMessage struct
			var notification NotificationMessage
			if err := json.Unmarshal(d.Body, &notification); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			// Store the notification in the database
			newNotification := models.Notification{
				UserID:    notification.UserID,
				Message:   notification.Message,
				CreatedAt: time.Now(),
				Read:      false,
			}
			db.DB.Create(&newNotification)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
