package main

import (
	"github.com/eniehack/simple-sns-go/app"
	"github.com/goadesign/goa"
)

// AuthorizationController implements the Authorization resource.
type AuthorizationController struct {
	*goa.Controller
}

// NewAuthorizationController creates a Authorization controller.
func NewAuthorizationController(service *goa.Service) *AuthorizationController {
	return &AuthorizationController{Controller: service.NewController("AuthorizationController")}
}

// Login runs the login action.
func (c *AuthorizationController) Login(ctx *app.LoginAuthorizationContext) error {
	// AuthorizationController_Login: start_implement

	// Put your logic here

	res := &app.Login{}
	return ctx.OK(res)
	// AuthorizationController_Login: end_implement
}

// Register runs the register action.
func (c *AuthorizationController) Register(ctx *app.RegisterAuthorizationContext) error {
	// AuthorizationController_Register: start_implement

	// Put your logic here

	return nil
	// AuthorizationController_Register: end_implement
}
