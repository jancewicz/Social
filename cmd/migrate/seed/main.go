package main

import (
	"log"
	"os"

	"github.com/jancewicz/social/internal/db"
	"github.com/jancewicz/social/internal/store"
)

func main() {
	addr := os.Getenv("DB_ADDR")

	connection, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()

	store := store.NewStorage(connection)

	db.Seed(store)
}
