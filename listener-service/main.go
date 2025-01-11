package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/BlackSound1/go-microservices/listener/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	// Try to connect to RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// Start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// Create a consumer to consume messages from queue
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// Watch queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

// connect retries connecting to RabbitMQ up to 5 times, with an exponentially
// increasing backoff time if there is an error. If the 5th attempt fails, an
// error is returned. Otherwise, the connection is returned.
func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	// Don't continue until RabbitMQ is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq-service") // Try to connect

		// If there's an error, add to the number of counts
		if err != nil {
			fmt.Println("RabbitMQ not ready yet...")
			counts++
		} else {
			// Otherwise, we have the connection
			connection = c
			log.Println("Connected to RabbitMQ")
			break
		}

		// If we've tried 5 times, give up
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		// Wait and try again. The time we wait is the number of counts squared,
		// leading to an exponentially increasing backoff time
		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backoff)
	}

	return connection, nil
}
