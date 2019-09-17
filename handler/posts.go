package handler

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid"
)

func (h *Handler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	requestData := new(CreatePostParams)

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		_, err := token.Method.(*jwt.SigningMethodRSA)
		if !err {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return LoadPrivateKey()
		}
	})

	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(MakeErrorResponseBody(http.StatusUnauthorized, "invaild authorization token"))
		return
	}

	// Bind Request by requestdata.

	if err := json.NewDecoder(r.Body).Decode(requestData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(MakeErrorResponseBody(http.StatusInternalServerError, "Internal Server Error. Please contact Admin."))
		return
	}

	if err := h.validate.Struct(requestData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(MakeErrorResponseBody(http.StatusBadRequest, "invaild request format"))
		return
	}

	if err := h.InsertPost(token.Claims.(jwt.MapClaims), requestData); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(MakeErrorResponseBody(http.StatusInternalServerError, "Internal Server Error. Please contact Admin."))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (h *Handler) InsertPost(Claims jwt.MapClaims, requestData *CreatePostParams) error {

	db := h.DB
	UserID := Claims["aud"]
	ulid := ulid.MustNew(ulid.Now(), rand.Reader)

	BindParams := map[string]interface{}{
		"PostID": ulid.String(),
		"UserID": UserID,
		"Body":   requestData.Body,
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
