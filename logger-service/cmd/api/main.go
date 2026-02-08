package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var webPort = ":6001"
var mongoUrl string
var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	mongoClinet, err := ConnectMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClinet

	fmt.Println("MongoDB connected.")
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
		fmt.Println("MongoDB disconnected.")
	}()

	// Server
	app := Config{
		Models: *data.New(client),
	}
	srv := &http.Server{
		Addr:    webPort,
		Handler: app.routes(),
	}

	// Register rpc server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Println("Failed to register rpc server")
	}

	// start rpc server
	go app.rpcListen()

	// http server
	err = srv.ListenAndServe()
	log.Fatal(err)

}

// RPC Listen
const rpcPort = "5001"

func (app *Config) rpcListen() error {
	fmt.Println("Starting rpc server on port:" + rpcPort)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listener.Close()

	// loop forever
	for {
		rpcConn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}

}

func ConnectMongo() (*mongo.Client, error) {
	mongoUrl := os.Getenv("MONGODB_URI")

	var client *mongo.Client
	var err error

	maxAttempts := 10

	for i := 1; i <= maxAttempts; i++ {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		client, err = mongo.Connect(options.Client().ApplyURI(mongoUrl))
		if err == nil {
			err = client.Ping(ctx, nil)
		}

		cancel()

		if err == nil {
			fmt.Println("Connected to MongoDB")
			return client, nil
		}

		fmt.Printf("Mongo not ready... retrying (%d/%d)\n", i, maxAttempts)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to mongo after %d attempts: %v", maxAttempts, err)
}
