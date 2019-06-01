package design

import (
	. "goa.design/goa/v3/dsl"
)

var UserIDAttribute = func() {
	Description("User ID")
	Example("hogehoge")
	MaxLength(15)
	MinLength(1)
	Pattern(`[^a-zA-Z0-9_]+`)
}

var ScreenNameAttribute = func() {
	Example("ほげほげ")
	MaxLength(20)
}

var PasswordAttribute = func() {
	Example("testpassword")
}

var BodyAttribute = func() {
	Example("にゃーん")
}

var NewPostPayload = Type("NewPostPayload", func() {
	Attribute("body", String, "")
	Token("token", String, "JWT Token.")
	Required("body", "token")
})

var DeletePostPayload = Type("DeletePostPayload", func() {
	Attribute("post_id", String, "unique id of the post.")
	Token("token", String, "JWT Token.")
	Required("post_id", "token")
})

var NewAccountPayload = Type("NewAccountPayload", func() {
	Attribute("userid", String, "unique id of the user.", UserIDAttribute)
	Attribute("screen_name", String, "screen name of the user.", ScreenNameAttribute)
	Attribute("password", String, "password of the user.", PasswordAttribute)
	Required("userid", "screen_name", "password")
})

var LoginPayload = Type("LoginPayload", func() {
	Attribute("userid", String, "", UserIDAttribute)
	Attribute("password", String, "", PasswordAttribute)
	Required("userid", "password")
})
