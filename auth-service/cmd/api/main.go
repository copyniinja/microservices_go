package main

import (
	"auth-service/data"
	"auth-service/infra/db"
	"fmt"
	"log"
	"net/http"
)

const webPort = 5000

type Config struct {
	Models data.Models
}

func main() {

	conn := db.NewConnection()

	model := data.New(conn)

	app := Config{
		Models: *model,
	}
	log.Println("Authentication service is ready.")

	srv := NewServer(app.routes())

	srv.Start()

}

type Server struct {
	server *http.Server
}

func NewServer(hnd http.Handler) *Server {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: hnd,
	}

	return &Server{
		server: server,
	}
}

func (s *Server) Start() {
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
