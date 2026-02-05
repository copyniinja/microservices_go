package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (app *Config) sendEmail(w http.ResponseWriter, r *http.Request) {

	var req struct {
		To          string         `json:"to"`
		Subject     string         `json:"subject"`
		Template    string         `json:"template"`
		Data        map[string]any `json:"data"`
		Attachments []string       `json:"attachments"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	msg := Message{
		To:          req.To,
		Subject:     req.Subject,
		DataMap:     req.Data,
		Attachments: req.Attachments,
	}

	go func() {
		if err := app.Mailer.Send("./templates/"+req.Template, msg); err != nil {
			log.Println("failed to send email:", err)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("email is being sent"))
}
