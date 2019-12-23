package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eniehack/persona-server/config"
	"github.com/eniehack/persona-server/handler"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"gopkg.in/go-playground/validator.v9"
)

func SetUpDataBase(Config *config.DatabaseConfig) (*sqlx.DB, error) {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", Config.User, Config.Password, Config.Host, Config.Database, Config.SSL)
	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect Database: %s", err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations/postgres",
	}
	_, err = migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		return nil, fmt.Errorf("failed migrations: %s", err)
	} /*else {
			log.Println("Applied %d migrations", n)
	} */
	return db, nil
}

func main() {
	var configFilePath string

	flag.StringVar(&configFilePath, "config", "./configs/config.toml", "config file's path.")
	flag.Parse()

	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load config file: %s", err)
	}

	corsSettings := cors.New(cors.Options{
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodDelete},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		MaxAge:         3600,
		ExposedHeaders: []string{"Authorization"},
	})

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsSettings.Handler)

	log.Println("Persona v0.1.0-alpha.1 starting......")

	db, err := SetUpDataBase(config)
	if err != nil {
		log.Fatalf(err)
	}

	db.SetConnMaxLifetime(1)
	defer db.Close()

	log.Println("finiched set up database")

	validator := validator.New()

	h := &handler.Handler{
		DB:       db,
		Validate: validator,
	}

	r.Mount("/api/v1", Router(h))
	log.Println("Finished to mount router.")

	if os.Getenv("PORT") == "" {
		log.Fatal(http.ListenAndServe(":3000", r))
	} else {
		log.Fatal(
			http.ListenAndServe(
				fmt.Sprintf(":%s", os.Getenv("PORT")),
				r,
			),
		)
	}
}
