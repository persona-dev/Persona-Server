// Code generated by goa v3.0.2, DO NOT EDIT.
//
// Authorization endpoints
//
// Command:
// $ goa gen github.com/eniehack/persona-server/design

package authorization

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "Authorization" service endpoints.
type Endpoints struct {
	Login    goa.Endpoint
	Register goa.Endpoint
}

// NewEndpoints wraps the methods of the "Authorization" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Login:    NewLoginEndpoint(s),
		Register: NewRegisterEndpoint(s),
	}
}

// Use applies the given middleware to all the "Authorization" service
// endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Login = m(e.Login)
	e.Register = m(e.Register)
}

// NewLoginEndpoint returns an endpoint function that calls the method "login"
// of service "Authorization".
func NewLoginEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*LoginPayload)
		return nil, s.Login(ctx, p)
	}
}

// NewRegisterEndpoint returns an endpoint function that calls the method
// "register" of service "Authorization".
func NewRegisterEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*NewAccountPayload)
		return nil, s.Register(ctx, p)
	}
}