package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// declareExchange declares a topic exchange for logs
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // Name
		"topic",      // Type
		true,         // Durable?
		false,        // Auto-deleted?
		false,        // Internal?
		false,        // No-wait?
		nil,          // Arguments
	)
}

// declareRandomQueue declares a random queue
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // Name
		false, // Durable?
		false, // Auto-deleted?
		true,  // Exclusive?
		false, // No-wait?
		nil,   // Arguments
	)
}
