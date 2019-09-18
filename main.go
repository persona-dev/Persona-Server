package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eniehack/persona-server/handler"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
	"gopkg.in/go-playground/validator.v9"
)

func SetUpDataBase(DataBaseName string) (*sqlx.DB, error) {
	switch DataBaseName {
	case "sqlite3":
		db, err := sqlx.Open("sqlite3", "test.db")
		if err != nil {
			return nil, fmt.Errorf("failed to connect Database: %s", err)
		}

		migrations := &migrate.FileMigrationSource{
			Dir: "migrations/sqlite3",
		}
		_, err = migrate.Exec(db.DB, "sqlite3", migrations, migrate.Up)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("failed migrations: %s", err)
		} /* else {
			log.Println("Applied %d migrations", n)
		} */
		return db, nil
	case "postgres":
		db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
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
	default:
		return nil, fmt.Errorf("invaild database flag")
	}
}

func main() {
	var DataBaseName string

	flag.StringVar(&DataBaseName, "database", "sqlite3", "Database name. sqlite3 or postgres.")
	flag.Parse()

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

	db, err := SetUpDataBase(DataBaseName)
	if err != nil {
		log.Fatalln(err)
	}

	db.SetConnMaxLifetime(1)
	defer db.Close()

	log.Println("finiched set up database")

	validator := validator.New()

	h := &handler.Handler{
		DB:       db,
		Validate: validator,
	}

	r.Mount("/api/v1", MainRouter(h))
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
