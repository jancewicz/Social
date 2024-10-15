package main

import (
	"log"
	"os"

	"github.com/jancewicz/social/internal/db"
	"github.com/jancewicz/social/internal/env"
	"github.com/jancewicz/social/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString(os.Getenv("SRV_ADDR"), ":8080"),
		db: dbConfig{
			addr:         env.GetString(os.Getenv("DB_ADDR"), "postgres://user:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt(os.Getenv("DB_MAX_OPEN_CONNS"), 30),
			maxIdleConns: env.GetInt(os.Getenv("DB_MAX_IDLE_CONNS"), 30),
			maxIdleTime:  env.GetString(os.Getenv("DB_MAX_IDLE_TIME"), "15min"),
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
