package main

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (a *Config) Broker(w http.ResponseWriter, r *http.Request) {

	resp := &JsonResponse{
		Success: true,
		Message: "Welcome to broker service",
	}
	json.NewEncoder(w).Encode(resp)
}
