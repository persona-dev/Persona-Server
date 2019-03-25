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

// Register runs the register action.
func (c *AuthorizationController) Register(ctx *app.RegisterAuthorizationContext) error {
	// AuthorizationController_Register: start_implement

	// Put your logic here

	return nil
	// AuthorizationController_Register: end_implement
}
