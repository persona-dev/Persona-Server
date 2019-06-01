// Code generated by goa v3.0.2, DO NOT EDIT.
//
// Post endpoints
//
// Command:
// $ goa gen github.com/eniehack/persona-server/design

package post

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "Post" service endpoints.
type Endpoints struct {
	Create    goa.Endpoint
	Reference goa.Endpoint
	Delete    goa.Endpoint
}

// NewEndpoints wraps the methods of the "Post" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		Create:    NewCreateEndpoint(s, a.JWTAuth),
		Reference: NewReferenceEndpoint(s),
		Delete:    NewDeleteEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "Post" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Create = m(e.Create)
	e.Reference = m(e.Reference)
	e.Delete = m(e.Delete)
}

// NewCreateEndpoint returns an endpoint function that calls the method
// "create" of service "Post".
func NewCreateEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*NewPostPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "JWT",
			Scopes:         []string{"api:access", "api:admin"},
			RequiredScopes: []string{"api:access"},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.Create(ctx, p)
	}
}

// NewReferenceEndpoint returns an endpoint function that calls the method
// "reference" of service "Post".
func NewReferenceEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*Post)
		return nil, s.Reference(ctx, p)
	}
}

// NewDeleteEndpoint returns an endpoint function that calls the method
// "delete" of service "Post".
func NewDeleteEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DeletePostPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "JWT",
			Scopes:         []string{"api:access", "api:admin"},
			RequiredScopes: []string{"api:access"},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.Delete(ctx, p)
	}
}