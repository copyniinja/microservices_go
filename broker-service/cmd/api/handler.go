package main

import (
	"broker-service/cmd/clients"
	"context"
	"encoding/json"
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
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Auth service error: " + err.Error()))
		return
	}

	//  Respond back to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResp)
}
