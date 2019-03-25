package main

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/eniehack/simple-sns-go/app"
	"github.com/goadesign/goa"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid"
)

// PostController implements the Post resource.
type PostController struct {
	*goa.Controller
}

// NewPostController creates a Post controller.
func NewPostController(service *goa.Service) *PostController {
	return &PostController{Controller: service.NewController("PostController")}
}

// Create runs the create action.
func (c *PostController) Create(ctx *app.CreatePostContext) error {
	// PostController_Create: start_implement

	// Put your logic here

	// jwtバリデーション

	// XSS対策にtext/templateを導入？

	// ULID生成

	now := time.Now()
	timestamp := now.Unix()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(now.UnixNano())), 0)
	ulid := ulid.MustNew(ulid.Timestamp(timestamp), entropy)

	db, err := sql.Open("sqlite", "test.db")

	if err != nil {
		log.Error(err)
		return ctx.InternalServerError(err)
	}

	if _, err := db.Exec(
		"INSERT INTO posts (post_id, user_id, body, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		ulid,
		"testuser",
		ctx.Body,
		now,
		now,
	); err != nil {
		log.Error(err)
		return ctx.InternalServerError(err)
	}

	return nil
	// PostController_Create: end_implement
}

// Delete runs the delete action.
func (c *PostController) Delete(ctx *app.DeletePostContext) error {
	// PostController_Delete: start_implement

	// Put your logic here

	return nil
	// PostController_Delete: end_implement
}

// Reference runs the reference action.
func (c *PostController) Reference(ctx *app.ReferencePostContext) error {
	// PostController_Reference: start_implement

	// Put your logic here

	var (
		createdAt  time.Time
		userID     string
		screenName string
		body       string
	)

	// PostIDが有効なものであるか確認する
	if _, err := ulid.Parse(ctx.PostID); err != nil {
		log.Error(err)
		return ctx.BadRequest(err)
	}

	db, err := sql.Open("sqlite", "test.db")
	if err != nil {
		log.Error(err)
	}

	defer db.Close()

	rows, err := db.Query(
		"SELECT posts.user_id, users.screen_name, posts.body, posts.created_at FROM posts JOIN users ON posts.user_id = users.user_id WHERE post_id = ?",
		ctx.PostID,
	)

	if err != nil {
		log.Error(err)
	}

	if err := rows.Scan(&userID, &screenName, &body, &createdAt); err != nil {
		log.Error(err)
	}

	res := &app.Post{}

	res.PostedAt = createdAt
	res.ScreenName = screenName
	res.Body = body
	res.UserID = userID

	return ctx.OK(res)
	// PostController_Reference: end_implement
}
