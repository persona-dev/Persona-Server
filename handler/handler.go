package handler

import (
	"database/sql"
	"errors"
)

type (
	Handler struct {
		DB *sql.DB
	}
	Argon2Params struct {
		memory      uint32
		iterations  uint32
		parallelism uint8
		saltLength  uint32
		keyLength   uint32
	}
	RegisterParams struct {
		UserID     string
		EMail      string
		ScreenName string
		Password   string
	}
)

var (
	ErrInvaildHash         = errors.New("the encoded hash is not in the correct format.")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)
