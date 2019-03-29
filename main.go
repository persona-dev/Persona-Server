package main

import (
	"database/sql"

	"github.com/eniehack/simple-sns-go/handler"

	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	pubkey, _ := ioutil.ReadFile("public-key.pub")
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	config := middleware.JWTConfig{
		SigningKey:    pubkey,
		SigningMethod: "RS512",
		Claims:        handler.Claims{},
	}

	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		e.Logger.Fatal(err)
	}

	h := &handler.Handler{DB: db}

	Authg := e.Group("/api/v1/auth")
	Authg.POST("/signature", h.Login)
	Authg.POST("/new", h.Register)

	Postg := e.Group("api/v1/posts")
	Postg.POST("/new", h.CreatePosts, middleware.JWTWithConfig(config))

	e.Logger.Fatal(e.Start(":8080"))
}
