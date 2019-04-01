package handler

import (
	"database/sql"
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type (
	Handler struct {
		DB *sql.DB
	}

	Claims struct {
		Scope string `json:"scope"`
		jwt.StandardClaims
	}

	Argon2Params struct {
		memory      uint32
		iterations  uint32
		parallelism uint8
		saltLength  uint32
		keyLength   uint32
	}
)

var (
	ErrInvaildHash         = errors.New("the encoded hash is not in the correct format.")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)
