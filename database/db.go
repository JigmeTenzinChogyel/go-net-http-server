package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func NewPostgresStorage() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	conStr := os.Getenv("DB_CONN_STR")

	if conStr == "" {
		conStr = "postgres://postgres:admin@localhost:5432/ecommerce?sslmode=disable"
	}
	db, err := sql.Open("postgres", "postgres://postgres:admin@localhost:5432/ecommerce?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	initPostgresStorage(db)

	return db, nil
}

func initPostgresStorage(db *sql.DB) {
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}
	log.Println("DB successfully connected!")
}
