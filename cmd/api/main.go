package main

import (
	"log"
	"os"
	"time"

	"github.com/jancewicz/social/internal/db"
	"github.com/jancewicz/social/internal/env"
	"github.com/jancewicz/social/internal/mailer"
	"github.com/jancewicz/social/internal/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
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
		addr:        env.GetString(os.Getenv("SRV_ADDR"), ":8080"),
		apiURL:      env.GetString(os.Getenv("EXTERNAL_URL"), "localhost:8080"),
		frontendURL: env.GetString(os.Getenv("FRONTEND_URL"), "http://localhost/4000"),
		db: dbConfig{
			addr:         os.Getenv("DB_ADDR"),
			maxOpenConns: env.GetInt(os.Getenv("DB_MAX_OPEN_CONNS"), 30),
			maxIdleConns: env.GetInt(os.Getenv("DB_MAX_IDLE_CONNS"), 30),
			maxIdleTime:  env.GetString(os.Getenv("DB_MAX_IDLE_TIME"), "15m"),
		},
		env: os.Getenv("ENV"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days to accpet invite
			fromEmail: os.Getenv("FROM_EMAIL"),
			mailTrap: mailTrapConfig{
				apiKey: os.Getenv("MAILTRAP_API_KEY"),
			},
		},
	}
	log.Println("DB Address:", cfg.db.addr)

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection established")

	store := store.NewStorage(db)

	// Mailer
	// mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	mailtrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailtrap,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
