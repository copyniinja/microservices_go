package main

import (
	"broker-service/cmd/clients"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const webPort = 4000

type Config struct {
	clients *clients.Clients
}

func main() {

	authBaseUrl := "http://auth-service:4000"
	app := Config{
		clients: clients.NewClients(authBaseUrl),
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
