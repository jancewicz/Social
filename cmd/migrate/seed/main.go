package main

import (
	"log"
	"os"

	"github.com/jancewicz/social/internal/db"
	"github.com/jancewicz/social/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("DB_ADDR")

	connection, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()

	store := store.NewStorage(connection)

	db.Seed(store, connection)
}
