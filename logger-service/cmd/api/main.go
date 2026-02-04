package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var webPort = ":6000"
var mongoUrl string
var client *mongo.Client

type Config struct {
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
}

func ConnectMongo() (*mongo.Client, error) {
	mongoUrl = os.Getenv("MONGODB_URI")

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(mongoUrl))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("could not ping mongodb:%v", err)
	}

	return mongoClient, nil
}
