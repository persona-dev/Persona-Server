package design

import (
	"os"

	. "goa.design/goa/v3/dsl"
)

var _ = API("Persona", func() {
	Title("Persona")
	Description("Layer構造を持つSNSのAPIです.")
	HTTP(func() {
		Consumes("application/json")
		Produces("application/json")
	})
	Server("persona", func() {
		Services("Authorization", "Post")
		Host("development", func() {
			URI("http://localhost:8000/api/v1")
		})

		Host("production", func() {
			URI("https://localhost:{port}/api/v1")
			Variable("port", String, "Port number", func() {
				Default(os.Getenv("PORT"))
			})
		})
	})
	Version("0.1.0")
})
