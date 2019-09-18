package main

import (
	"log"
	"net/http"

	"github.com/eniehack/persona-server/handler"
	"github.com/eniehack/persona-server/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

func DefineCORS() {}

func MainRouter(h *handler.Handler) http.Handler {

	rsapublickey, err := utils.ReadPublicKey()
	if err != nil {
		log.Fatalf("init(): Failed to road RSA public key: %s", err.Error())
	}
	tokenAuth := jwtauth.New("RS512", rsapublickey, nil)

	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signature", h.Login)
		r.Post("/new", h.Register)
	})

	r.Route("/posts", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/new", h.CreatePosts)
	})

	return r
}
