package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5"
)

const DATABASE_URL = `postgres://username:password@localhost:5432/database_name`

func NewConnection() *sql.DB {
	conn, err := sql.Open("pgx", DATABASE_URL)

	if err != nil {
		log.Fatalf("Unable to connect to database:%v", err)

	}

	return conn
}
