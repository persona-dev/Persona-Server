package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var JWT = JWTSecurity("JWT", func() {
	TokenURL("/auth/signature")
	Header("Authorization")
	Scope("api:access", "API access")
	Scope("api:admin", "管理者によるアクセス権限")
})

var _ = Resource("Authorization", func() {
	BasePath("/auth")
	Security(JWT, func() {
		Scope("api:access")
	})
	Action("register", func() {
		Description("新規登録")
		Routing(
			POST("/new"),
		)
		Payload(NewAccountPayload, func() {
			Required("userid", "screen_name", "password")
		})
		Response(Created, "/users/[0-9a-z]+")
		Response(Conflict, ErrorMedia)
		UseTrait("error")
	})
})

var _ = Resource("Post", func() {
	BasePath("/posts")
	Security(JWT, func() {
		Scope("api:access")
	})
	Action("create", func() {
		Description("新規投稿")
		Routing(
			POST("/new"),
		)
		Payload(NewPostPayload, func() {
			Required("body")
		})
		Response(NoContent)
		UseTrait("error")
	})
	/*Action("timeline", func(){
		Description("タイムラインの更新")
		Routing(
			GET("/")
		)
		Response(OK)
		Response()
	})*/
	Action("reference", func() {
		Description("投稿の参照")
		NoSecurity()
		Routing(
			GET("/:post_id"),
		)
		Params(func() {
			Param("post_id", String, "投稿固有のID", func() {
				Pattern("[0-9A-Z]{26}")
			})
		})
		Response(OK, PostDetailsMedia)
		UseTrait("error")
	})
	Action("delete", func() {
		Description("投稿の削除")
		Routing(
			DELETE("/:post_id"),
		)
		Params(func() {
			Param("post_id", String, "投稿固有のID", func() {
				Pattern("[0-9A-Z]{26}")
			})
		})
		Response(NoContent)
		UseTrait("error")
	})
})
