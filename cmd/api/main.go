package main

import (
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jancewicz/social/internal/auth"
	"github.com/jancewicz/social/internal/db"
	"github.com/jancewicz/social/internal/env"
	"github.com/jancewicz/social/internal/mailer"
	"github.com/jancewicz/social/internal/store"
	"github.com/jancewicz/social/internal/store/cache"
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
		frontendURL: env.GetString(os.Getenv("FRONTEND_URL"), "http://localhost:5173/"),
		db: dbConfig{
			addr:         os.Getenv("DB_ADDR"),
			maxOpenConns: env.GetInt(os.Getenv("DB_MAX_OPEN_CONNS"), 30),
			maxIdleConns: env.GetInt(os.Getenv("DB_MAX_IDLE_CONNS"), 30),
			maxIdleTime:  env.GetString(os.Getenv("DB_MAX_IDLE_TIME"), "15m"),
		},
		redisCfg: redisConfig{
			addr:     os.Getenv("REDIS_ADDR"),
			password: os.Getenv("REDIS_PASSWORD"),
			db:       env.GetInt(os.Getenv("REDIS_DB"), 0),
			enable:   env.GetBool(os.Getenv("REDIS_ENABLED"), true),
		},
		env: os.Getenv("ENV"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days to accpet invite
			fromEmail: os.Getenv("FROM_EMAIL"),
			mailTrap: mailTrapConfig{
				apiKey: os.Getenv("MAILTRAP_API_KEY"),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user:     os.Getenv("AUTH_BASIC_USER"),
				password: os.Getenv("AUTH_BASIC_PASS"),
			},
			token: tokenConfig{
				secret: os.Getenv("AUTH_TOKEN_SECRET"),
				exp:    time.Hour * 24 * 3, // Three days
				issuer: os.Getenv("TOKEN_ISSUER"),
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

	// Redis cache
	var redisDB *redis.Client
	if cfg.redisCfg.enable {
		redisDB = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.password, cfg.redisCfg.db)
		logger.Info("redis connection established")
	}

	store := store.NewStorage(db)
	cacheStore := cache.NewRedisStorage(redisDB)

	// Mailer
	mailtrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	// Authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.issuer,
		cfg.auth.token.issuer,
	)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStore:    cacheStore,
		logger:        logger,
		mailer:        mailtrap,
		authenticator: jwtAuthenticator,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
