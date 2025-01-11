package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// NewConsumer creates a new Consumer by setting up the RabbitMQ connection and
// declaring the topic exchange.
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{conn: conn}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// setup sets up the RabbitMQ consumer by declaring the exchange.
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

func (consumer *Consumer) Listen(topics []string) error {

	// Try to get the channel
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Try to get the queue
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	// Go through each topic and bind it to the queue
	for _, s := range topics {
		err := ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer name
		true,   // Auto-ack?
		false,  // Exclusive?
		false,  // No local?
		false,  // No wait?
		nil,    // Args
	)
	if err != nil {
		return err
	}

	// Create a channel that runs in the background forever
	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload

			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)

		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

// handlePayload takes a payload and handles it in various ways depending on its type
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// Log whatever we get
		err := logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}
	case "auth":
		// Authenticate

	default:
		err := logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// logEvent sends the given payload entry to the log service as a JSON object.
// It performs an HTTP POST request to the log service URL and checks for a successful response.
func logEvent(entry Payload) error {
	// Convert given entry into JSON
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	// Create a new request to the log service
	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set the necessary header
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
