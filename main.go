//go:generate goagen bootstrap -d github.com/eniehack/simple-sns-go/design

package main

import (
	"github.com/eniehack/simple-sns-go/app"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

func main() {
	// Create service
	service := goa.New("Simple-SNS")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "Authorization" controller
	c := NewAuthorizationController(service)
	app.MountAuthorizationController(service, c)
	// Mount "Post" controller
	c2 := NewPostController(service)
	app.MountPostController(service, c2)

	// Start service
	if err := service.ListenAndServeTLS(":8080", "cert.pem", "key.pem"); err != nil {
		service.LogError("startup", "err", err)
	}

}
