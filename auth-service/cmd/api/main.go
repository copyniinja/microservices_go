package main

import (
	"auth-service/data"
	"database/sql"
	"log"
)

const webPort = 5000

type Config struct {
	Db     *sql.DB
	Models data.Models
}

func main() {

	log.Println("Authentication service is ready.")
}
