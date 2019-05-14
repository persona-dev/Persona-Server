package handler

import (
	"crypto/rand"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/oklog/ulid"
)

func (h *Handler) CreatePosts(c echo.Context) error {

	user := c.Get("token").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["aud"]
	//fmt.Println(userID)

	body := c.FormValue("body")
	//fmt.Println(body)
	if body == "" {
		log.Println("No body.")
		return echo.ErrInternalServerError
	}

	now := time.Now().Format(time.RFC3339Nano)
	ulid := ulid.MustNew(ulid.Now(), rand.Reader)

	db := h.DB

	if _, err := db.Exec(
		"INSERT INTO posts (post_id, user_id, body, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		ulid.String(),
		userID,
		body,
		now,
		now,
	); err != nil {
		log.Println("INSERT Err", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
