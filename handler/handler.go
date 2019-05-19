package handler

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type (
	Handler struct {
		DB *sqlx.DB
	}
	Argon2Params struct {
		memory      uint32
		iterations  uint32
		parallelism uint8
		saltLength  uint32
		keyLength   uint32
	}
	RegisterParams struct {
		UserID     string `json:"userid" validate:"required,max=15"`
		EMail      string `json:"email" validate:"required,email"`
		ScreenName string `json:"screen_name" validate:"required,max=50"`
		Password   string `json:"password" validate:"required"`
	}
	LoginParams struct {
		UserName string `json:"userid" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)

var (
	ErrInvaildHash         = errors.New("the encoded hash is not in the correct format.")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)
