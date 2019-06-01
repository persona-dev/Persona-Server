package design

import (
	. "goa.design/goa/v3/dsl"
)

var PostDetailsMedia = ResultType("application/vnd.post+json", func() {
	Description("投稿の詳細な情報を返す際のレスポンス")
	Attributes(func() {
		Attribute("posted_at", String, "投稿日時", func() {
			Format(FormatDateTime)
		})
		Attribute("user_id", String, "投稿者", UserIDAttribute)
		Attribute("screen_name", String, ScreenNameAttribute)
		Attribute("body", String, "投稿内容", BodyAttribute)
		Required("posted_at", "user_id", "screen_name", "body")
	})
})

var LoginMedia = ResultType("application/vnd.login+json", func() {
	Description("ログインが成功した際に認証トークンを返すレスポンス")
	Attributes(func() {
		Attribute("token", String, "認証トークン", func() {
			Example("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.cThIIoDvwdueQB468K5xDc5633seEFoqwxjF_xSJyQQ")
		})

		Required("token")
	})
})
