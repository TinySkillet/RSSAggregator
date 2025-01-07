package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// separate storage struct so that if we want,
// we can change postgres with smth else
type Storage struct {
	store *PostgresStore
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() *PostgresStore {
	connString := os.Getenv("DBCONN_STR")
	if connString == "" {
		log.Fatal("DBCONN_STR not found in environment!")
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Unable to connect to the database server! Error: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to communicate with the database server! Error: %v", err)
	}

	log.Print("Connected with postgres server succesfully!")
	return &PostgresStore{
		db: db,
	}
}
