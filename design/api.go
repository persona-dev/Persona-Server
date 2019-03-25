package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("Simple-SNS", func() {
	Title("Simple-SNS")
	Description("Layer構造を持つSNSのAPIです.")
	Version("0.1")
	Scheme("https")
	Host("localhost:8080")
	BasePath("/api/v1")
	Consumes("application/json")
	Produces("application/json")
	Trait("error", func() {
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
	ResponseTemplate(Created, func(pattern string) {
		Description("リソースの作成が完了しました。")
		Status(201)
		Headers(func() {
			Header("Location", func() {
				Pattern(pattern)
			})
		})
	})
})
