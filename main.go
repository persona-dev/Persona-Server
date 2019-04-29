package main

import (
	"crypto/rsa"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/eniehack/simple-sns-go/handler"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		e.Logger.Fatal("db connection", err)
	}
	db.SetConnMaxLifetime(1)
	defer db.Close()

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations/sqlite3",
	}
	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Applied %d migrations", n)
	}

	h := &handler.Handler{DB: db}

	Authg := e.Group("/api/v1/auth")
	Authg.POST("/signature", h.Login)
	Authg.POST("/new", h.Register)

	Postg := e.Group("/api/v1/posts")
	Postg.Use(JWTAuthentication)
	Postg.POST("/new", h.CreatePosts)

	e.Logger.Fatal(e.Start(":8080"))
}
