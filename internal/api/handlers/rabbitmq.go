
package handlers

import (
	"log"

	"github.com/streadway/amqp"
)

var (
	RabbitMQConn    *amqp.Connection
	RabbitMQChannel *amqp.Channel
)

func InitRabbitMQ() {
	log.Println("Connecting to RabbitMQ...")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	RabbitMQConn = conn
	log.Println("RabbitMQ connection established.")

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	RabbitMQChannel = channel
	log.Println("RabbitMQ channel opened.")

	_, err = channel.QueueDeclare(
		"imageProcessingQueue",
		true,  
		false, 
		false, 
		false, 
		nil,   
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
	log.Println("Queue 'imageProcessingQueue' declared successfully.")
}

func CloseRabbitMQ() {
	if RabbitMQChannel != nil {
		if err := RabbitMQChannel.Close(); err != nil {
			log.Println("Error closing RabbitMQ channel:", err)
		}
	}
	if RabbitMQConn != nil {
		if err := RabbitMQConn.Close(); err != nil {
			log.Println("Error closing RabbitMQ connection:", err)
		}
	}
	log.Println("RabbitMQ connection closed.")
}
