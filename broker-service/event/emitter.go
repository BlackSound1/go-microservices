package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

// Setup initializes the RabbitMQ channel and declares an exchange
func (e *Emitter) Setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	return declareExchange(channel)
}

// Push sends a message to the RabbitMQ server with the given event and severity
func (e *Emitter) Push(event string, severity string) error {

	// Try to get the channel
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel...")

	// Try to publish to the channel
	err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// NewEventEmitter returns an Emitter that can be used to send messages to RabbitMQ
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {

	emitter := Emitter{connection: conn}

	// Try to set up the emitter
	err := emitter.Setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
