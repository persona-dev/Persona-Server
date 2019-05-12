package main

import (
	"crypto/rsa"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/eniehack/simple-sns-go/handler"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
)

func JWTAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		AuthorizationHeader := c.Request().Header.Get("Authorization")
		SplitAuthorization := strings.Split(AuthorizationHeader, " ")
		if SplitAuthorization[0] != "Bearer" {
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "invalid token.",
			}
		}
		token, err := jwt.Parse(SplitAuthorization[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return LookupPublicKey()
		})
		if err != nil || !token.Valid {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  "invalid token.",
				Internal: err,
			}
		}
		c.Set("token", token)
		return next(c)
	}
}

func LookupPublicKey() (*rsa.PublicKey, error) {
	Key, err := ioutil.ReadFile("public-key.pem")
	if err != nil {
		return nil, err
	}
	ParsedKey, err := jwt.ParseRSAPublicKeyFromPEM(Key)
	return ParsedKey, err
}

func SetUpDataBase(DataBaseName string) (*sqlx.DB, error) {
	switch DataBaseName {
	case "sqlite3":
		db, err := sqlx.Open("sqlite3", "test.db")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to connect Database: %s", err))
		}

		migrations := &migrate.FileMigrationSource{
			Dir: "migrations/sqlite3",
		}
		_, err = migrate.Exec(db.DB, "sqlite3", migrations, migrate.Up)
		if err != nil {
			log.Println(err)
			return nil, errors.New(fmt.Sprintf("failed migrations: %s", err))
		} /* else {
			log.Println("Applied %d migrations", n)
		} */
		return db, nil
	case "postgres":
		db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to connect Database: %s", err))
		}

		migrations := &migrate.FileMigrationSource{
			Dir: "migrations/postgres",
		}
		_, err = migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed migrations: %s", err))
		} /*else {
				log.Println("Applied %d migrations", n)
		} */
		return db, nil
	default:
		return nil, errors.New("invaild database flag")
	}
}

func main() {
	var DataBaseName string

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods:  []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodDelete},
		AllowHeaders:  []string{"Authorization", "Content-Type"},
		MaxAge:        3600,
		ExposeHeaders: []string{"Authorization"},
	}))

	flag.StringVar(&DataBaseName, "database", "sqlite3", "Database name. sqlite3 or postgres.")
	flag.Parse()

	db, err := SetUpDataBase(DataBaseName)
	if err != nil {
		e.Logger.Fatal(err)
	}

	db.SetConnMaxLifetime(1)
	defer db.Close()

	h := &handler.Handler{DB: db}

	Authg := e.Group("/api/v1/auth")
	Authg.POST("/signature", h.Login)
	Authg.POST("/new", h.Register)

	Postg := e.Group("/api/v1/posts")
	Postg.Use(JWTAuthentication)
	Postg.POST("/new", h.CreatePosts)

	e.Logger.Fatal(
		e.Start(
			fmt.Sprintf(
				":%s", os.Getenv("PORT"),
			),
		),
	)
}
