package main

import (
	"github.com/eniehack/simple-sns-go/controller"

	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	pubkey, _ := ioutil.ReadFile("public-key.pub")
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	config := middleware.JWTConfig{
		SigningKey:    pubkey,
		SigningMethod: "RS512",
		Claims:        controller.Claims{},
	}

	Authg := e.Group("/api/v1/auth")
	Authg.POST("/signature", controller.Login)
	Authg.POST("/new", controller.Register)

	Postg := e.Group("api/v1/posts")
	Postg.POST("/new", controller.CreatePosts, middleware.JWTWithConfig(config))

	e.Logger.Fatal(e.Start(":8080"))
}
