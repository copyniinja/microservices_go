package main

import (
	"listener/event"
	"log"
	"math/rand"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const rabbitmqUrl = "amqp://guest:guest@rabbitmq:5762/"

func main() {
	// connect to rabbitmq
	conn := connectRabbitMQ(rabbitmqUrl)
	log.Println(conn)

	// start listening for messages
	log.Println("Listening for and consuming rabbitmq messages")

	// create consumers
	consumer, err := event.NewConsumer(conn)
	if err != nil {
		log.Fatal(err)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connectRabbitMQ(url string) *amqp.Connection {
	backoff := time.Second
	maxBackoff := 30 * time.Second

	for {
		conn, err := amqp.Dial(url)
		if err == nil {
			log.Println("Connected to rabbitmq.")
			return conn
		}

		// add jitter
		jitter := time.Duration(rand.Int63n(int64(time.Second)))
		wait := backoff + jitter

		log.Printf("RabbitMQ not ready. Retrying in %v...", wait)
		time.Sleep(wait)

		// Double backoff
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}

	}

}
