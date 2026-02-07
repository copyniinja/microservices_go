package main

import (
	"broker-service/cmd/clients"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = 4000

type Config struct {
	clients *clients.Clients
	Rabbit  *amqp.Connection
}

const rabbitmqUrl = "amqp://guest:guest@rabbitmq:5672/"

func main() {

	authBaseUrl := "http://auth-service:5000"
	loggerBaseUrl := "http://logger-service:6001"
	conn := connectRabbitMQ(rabbitmqUrl)
	app := Config{
		clients: clients.NewClients(authBaseUrl, loggerBaseUrl),
		Rabbit:  conn,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: app.routes(),
	}

	go func() {
		log.Println("The broker service is listening on port:", webPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}

	}()

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	fmt.Println("Shutdown signal received")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server shut down gracefully")

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
