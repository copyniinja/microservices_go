package main

import (
	"broker-service/cmd/clients"
	"broker-service/event"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type JsonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ReqBody struct {
	Message string `json:"message"`
}

func (a *Config) Broker(w http.ResponseWriter, r *http.Request) {
	var req ReqBody

	json.NewDecoder(r.Body).Decode(&req)
	log.Println("Req-Body:", req)
	w.Header().Set("Content-Type", "application/json")

	resp := &JsonResponse{
		Success: true,
		Message: "Welcome to broker service",
	}
	json.NewEncoder(w).Encode(resp)
}

type RequestPayload struct {
	Action string               `json:"action"`
	Auth   *clients.AuthPayload `json:"auth_payload,omitempty"`
	Log    *clients.LogPayload  `json:"log_payload,omitempty"`
}
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (a *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var reqPayload RequestPayload

	// Decode the json request body
	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	// Actions
	switch reqPayload.Action {
	case "auth":
		a.authenticate(w, *reqPayload.Auth)
	case "log":
		a.logEventViaRabbit(w, *reqPayload.Log)
	default:
		w.WriteHeader(400)
		w.Write([]byte("Unknown action"))
		return
	}

}

func (a *Config) authenticate(w http.ResponseWriter, payload clients.AuthPayload) {
	ctx := context.Background()
	// Call the auth microservices
	authResp, err := a.clients.Auth.Login(ctx, &payload)
	fmt.Println("ERROR:", err)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]interface{}{
			"error":   true,
			"message": "Auth service error: " + err.Error(),
		})
		return

	}

	//  Respond back to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResp)
}

func (a *Config) logItem(w http.ResponseWriter, payload clients.LogPayload) {
	ctx := context.Background()
	// Call the logger microservices
	logResp, err := a.clients.Log.Insert(ctx, &payload)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]interface{}{
			"error":   true,
			"message": "logger service error: " + err.Error(),
		})
		return
	}

	//  Respond back to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logResp)

}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	writeJSON(w, statusCode, payload)
	return nil
}
func (app *Config) logEventViaRabbit(w http.ResponseWriter, l clients.LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	writeJSON(w, http.StatusAccepted, payload)
}

// pushToQueue pushes a message into RabbitMQ
func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}
