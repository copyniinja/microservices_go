package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const DATABASE_URL = `postgres://faiyaz:faiyaz@localhost:5432/auth_db`

func NewConnection() *sql.DB {
	conn, err := sql.Open("pgx", DATABASE_URL)

	if err != nil {
		log.Fatalf("Unable to connect to database:%v", err)

	}

	if err := conn.Ping(); err != nil {
		log.Fatalf("DB not reachable: %v", err)
	}

	return conn
}
