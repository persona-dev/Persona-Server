package design

import (
	_ "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var NewPostPayload = Type("NewPostPayload", func() {
	Attribute("body", func() {
		Example("にゃーん")
		MaxLength(500)
	})
})

var NewAccountPayload = Type("NewAccountPayload", func() {
	Attribute("userid", func() {
		Example("hogehoge")
		MaxLength(15)
	})
	Attribute("screen_name", func() {
		Example("ほげほげ")
		MaxLength(20)
	})
	Attribute("password", func() {
		Example("testpassword")
	})
})
