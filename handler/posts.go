package handler

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/oklog/ulid"
)

func (h *Handler) CreatePosts(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*Claims)
	userID := claims.StandardClaims.Audience
	fmt.Println(userID)

	body := c.FormValue("body")
	fmt.Println(body)
	if body == "" {
		log.Println("No body.")
		return echo.ErrInternalServerError
	}

	now := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(now.UnixNano())), 0)
	ulid := ulid.MustNew(ulid.Timestamp(now), entropy)

	db := h.DB
	defer db.Close()

	if _, err := db.Exec(
		"INSERT INTO posts (post_id, user_id, body, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		ulid,
		userID,
		body,
		now,
		now,
	); err != nil {
		log.Println("INSERT Err", err)
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusOK)
}
