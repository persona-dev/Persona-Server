package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var PostDetailsMedia = MediaType("application/vnd.post+json", func() {
	Description("投稿の詳細な情報を返す際のレスポンス")
	Attributes(func() {
		Attribute("posted_at", DateTime, "投稿日時")
		Attribute("user_id", String, "投稿者", func() {
			Example("hogehoge")
			MaxLength(15)
		})
		Attribute("screen_name", String, "投稿者のスクリーンネーム", func() {
			Example("ほげほげ")
			MaxLength(20)
		})
		Attribute("body", String, "投稿内容", func() {
			Example("にゃーん")
		})

		Required("posted_at", "user_id", "screen_name", "body")
	})

	View("default", func() {
		Attribute("posted_at")
		Attribute("user_id")
		Attribute("screen_name")
		Attribute("body")
	})
})
