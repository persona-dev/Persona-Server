package handler

import (
	"crypto/rand"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
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

	ulid := ulid.MustNew(ulid.Now(), rand.Reader)

	db := h.DB

	BindParams := map[string]interface{}{
		"PostID": ulid.String(),
		"UserID": UserID,
		"Body":   body,
		"Now":    time.Now().Format(time.RFC3339Nano),
	}

	Query, Params, err := sqlx.Named(
		"INSERT INTO posts (post_id, user_id, body, created_at, updated_at) VALUES (:PostID, :UserID, :Body, :Now, :Now)",
		BindParams,
	)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status_code": "500",
		})
	}

	Rebind := db.Rebind(Query)

	if _, err := db.Exec(Rebind, Params); err != nil {
		log.Println("INSERT Err", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
