package event

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	loggerServiceURL = "http://logger-service:8080/log" // ← change this to your real logger URL
	httpTimeout      = 5 * time.Second
)

type Consumer struct {
	conn *amqp.Connection

	httpClient *http.Client
}

func NewConsumer(conn *amqp.Connection) (*Consumer, error) {
	httpClient := &http.Client{
		Timeout: httpTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	consumer := &Consumer{
		conn: conn,

		httpClient: httpClient,
	}

	if err := consumer.setup(); err != nil {
		return nil, err
	}

	return consumer, nil
}

func (c *Consumer) setup() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return declareExchange(ch)
}

// Payload is what we expect from the queue messages
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type LogResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		if err = ch.QueueBind(
			q.Name,       // queue
			topic,        // routing key = topic
			"logs_topic", // exchange
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to bind queue to topic %s: %w", topic, err)
		}
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer tag (empty = auto)
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	forever := make(chan struct{})

	go func() {
		for msg := range msgs {
			var payload Payload
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				log.Printf("failed to unmarshal message: %v → body: %q", err, msg.Body)
				continue
			}

			// handle in goroutine so one slow log doesn't block others
			go c.handlePayload(payload)
		}
	}()

	log.Printf("[*] Waiting for message [exchange: logs_topic, queue: %s, topics: %v]", q.Name, topics)
	<-forever

	return nil
}

func (c *Consumer) handlePayload(payload Payload) {
	switch payload.Name {

	case "log", "event":
		err := c.logEvent(payload)
		if err != nil {
			log.Println("logger error:", err)
		}

	case "auth":
		// Authentication

	default:
		err := c.logEvent(payload)
		if err != nil {
			log.Println("logger error:", err)
		}
	}
}

func (c *Consumer) logEvent(payload Payload) error {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loggerServiceURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("logger service returned status %d", resp.StatusCode)
	}

	var logResp LogResponse
	if err := json.NewDecoder(resp.Body).Decode(&logResp); err != nil {
		return fmt.Errorf("failed to decode logger response: %w", err)
	}

	if logResp.Error {
		return fmt.Errorf("logger service error: %s", logResp.Message)
	}

	return nil
}
