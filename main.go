package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	"github.com/eniehack/simple-sns-go/handler"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	pubkey, err := ioutil.ReadFile("public-key.pem")
	if err != nil {
		log.Println(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		e.Logger.Fatal("db connection", err)
	}

	h := &handler.Handler{DB: db}

	Authg := e.Group("/api/v1/auth")
	Authg.POST("/signature", h.Login)
	Authg.POST("/new", h.Register)

	Postg := e.Group("/api/v1/posts")
	config := middleware.JWTConfig{
		SigningKey:    pubkey,
		SigningMethod: "RS512",
		Claims:        &handler.Claims{},
	}
	Postg.Use(middleware.JWTWithConfig(config))
	Postg.POST("/new", h.CreatePosts)

	e.Logger.Fatal(e.Start(":8080"))
}
