package main

import (
	"log"
	"os"
	"time"

	"github.com/jancewicz/social/internal/db"
	"github.com/jancewicz/social/internal/env"
	"github.com/jancewicz/social/internal/store"
	"github.com/joho/godotenv"
)

const version = "0.0.2"

//	@title			Social API
//	@description	API for Social web app
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securitydefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr:   env.GetString(os.Getenv("SRV_ADDR"), ":8080"),
		apiURL: env.GetString(os.Getenv("EXTERNAL_URL"), "localhost:8080"),
		db: dbConfig{
			addr:         os.Getenv("DB_ADDR"),
			maxOpenConns: env.GetInt(os.Getenv("DB_MAX_OPEN_CONNS"), 30),
			maxIdleConns: env.GetInt(os.Getenv("DB_MAX_IDLE_CONNS"), 30),
			maxIdleTime:  env.GetString(os.Getenv("DB_MAX_IDLE_TIME"), "15m"),
		},
		env: os.Getenv("ENV"),
		mail: mailConfig{
			exp: time.Hour * 24 * 3, // 3 days to accpet invite
		},
	}
	log.Println("DB Address:", cfg.db.addr)

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("database connection established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
