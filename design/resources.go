package design

import (
	. "goa.design/goa/v3/dsl"
)

var JWT = JWTSecurity("JWT", func() {
	Scope("api:access", "API access")
	Scope("api:admin", "管理者によるアクセス権限")
})

var _ = Service("Authorization", func() {
	HTTP(func() {
		Path("/api/v1/auth")
	})
	Method("login", func() {
		Description("ログイン")
		Payload(LoginPayload)
		HTTP(func() {
			POST("/signature")
			Response(StatusOK)
		})
		/*
			Response(Unauthorized)
			UseTrait("error")
		*/
	})
	Method("register", func() {
		Description("新規登録")
		Payload(NewAccountPayload)
		HTTP(func() {
			POST("/new")
			Response(StatusCreated)
		})
		/*
			Response(Created, "/users/[0-9a-z]+")
			Response(Conflict, ErrorMedia)
			UseTrait("error")
		*/
	})
})

var _ = Service("Post", func() {
	HTTP(func() {
		Path("/api/v1/posts")
	})
	Security(JWT, func() {
		Scope("api:access")
	})
	Method("create", func() {
		Description("新規投稿")
		Security(JWT, func() {
			Scope("api:access")
		})
		Payload(NewPostPayload)
		HTTP(func() {
			POST("/new")
			Response(StatusNoContent)
		})
		/*
			UseTrait("error")
		*/
	})
	/*Action("timeline", func(){
		Description("タイムラインの更新")
		Routing(
			GET("/")
		)
		Response(OK)
		Response()
	})*/
	Method("reference", func() {
		Description("投稿の参照")
		NoSecurity()
		Payload(PostDetailsMedia)

		HTTP(func() {
			GET("/{post_id}")
			Params(func() {
				Param("post_id", String, "投稿固有のID", func() {
					Pattern("[0-9A-Z]{26}")
				})
				Required("post_id")
			})
			Response(StatusOK)
		})
		//PostDetailsMedia)
		//UseTrait("error")
	})
	Method("delete", func() {
		Description("投稿の削除")
		Payload(DeletePostPayload)
		HTTP(func() {
			DELETE("/{post_id}")
			Params(func() {
				Param("post_id", String, "投稿固有のID", func() {
					Pattern("[0-9A-Z]{26}")
				})
				Required("post_id")
			})
			Response(StatusNoContent)
		})
		/*
			UseTrait("error")
		*/
	})
})
