package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	Mailer *Mailer
}

var webPort = 6002

func main() {
	mailer := NewMailer(Mail{
		Domain:      "localhost",
		Port:        1025,
		Host:        "mailhog",
		FromAddress: "lol@gmail.com",
		FromName:    "Ajib",
	})
	app := &Config{
		Mailer: mailer,
	}
	err := mailer.Send(
		"./templates/welcome.html",
		Message{
			To:      "user@test.com",
			Subject: "Welcome!",
			DataMap: map[string]any{
				"Name": "Faiyaz",
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ email sent")
	log.Println(app)
	server := NewServer(webPort, app.routes())
	server.Start()
}

type Server struct {
	Port    string
	Handler http.Handler
}

func NewServer(port int, handler http.Handler) *Server {

	return &Server{
		Port:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
}

func (s *Server) Start() {
	srv := &http.Server{
		Addr:    s.Port,
		Handler: s.Handler,
	}
	log.Println("Mail server is ready.")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("failed to start the server:", err)
	}
}
