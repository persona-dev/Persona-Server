package controller

import (
	"database/sql"
)

type (
	Handler struct {
		DB *sql.DB
	}
)