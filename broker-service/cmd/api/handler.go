package main

import (
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
