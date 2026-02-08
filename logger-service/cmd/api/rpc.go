package main

import (
	"context"
	"fmt"
	"logger-service/data"
	"time"
)

type RPCServer struct {
}

type RPCPayload struct {
	Name string
	Data string
}

// function name must starts with a capital letter ( "___.LogInfo") and resp bust be a pointer to string
func (rpc *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	// Insert log into mongo
	collections := client.Database("logs").Collection("logs")
	_, err := collections.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	// Send response back to the client that called this function
	*resp = fmt.Sprintf("Processed payload via rpc: %s", payload.Name)
	return nil
}
