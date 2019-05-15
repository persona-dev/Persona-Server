package handler

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/oklog/ulid"
)

func (h *Handler) CreatePosts(c echo.Context) error {
	User := c.Get("token").(*jwt.Token)
	Claims := User.Claims.(jwt.MapClaims)

	Body := c.FormValue("body")
	if Body == "" {
		log.Println("No body.")
		return echo.ErrInternalServerError
	}

	if err := h.InsertPost(Claims, Body); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status_code": "500",
		})
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) InsertPost(Claims jwt.MapClaims, Body string) error {

	db := h.DB
	UserID := Claims["aud"]
	ulid := ulid.MustNew(ulid.Now(), rand.Reader)

	BindParams := map[string]interface{}{
		"PostID": ulid.String(),
		"UserID": UserID,
		"Body":   Body,
		"Now":    time.Now().Format(time.RFC3339Nano),
	}

	Query, Params, err := sqlx.Named(
		"INSERT INTO posts (post_id, user_id, body, created_at, updated_at) VALUES (:PostID, :UserID, :Body, :Now, :Now)",
		BindParams,
	)
	if err != nil {
		return fmt.Errorf("Error InsertPost(): failed to bind Parameters. %s", err)
	}

	Rebind := db.Rebind(Query)

	if _, err := db.Exec(Rebind, Params...); err != nil {
		return fmt.Errorf("Error InsertPost(): failed to insert user data. %s", err)
	}

	return nil
}
