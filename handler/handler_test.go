package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo"
)

/*
func TestMain(m *testing.M) {
	h, err := SetUpDataBase()
	db := h.DB
	ExitCode := m.Run()
	db.Close()
	os.Exit(ExitCode)
}
*/

func SetUpDataBase() (*Handler, error) {
	db, err := sqlx.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}

	h := &Handler{DB: db}

	migrations := &migrate.FileMigrationSource{
		Dir: "../migrations/sqlite3",
	}
	if _, err := migrate.Exec(db.DB, "sqlite3", migrations, migrate.Up); err != nil {
		return nil, err
	}
	return h, err
}

func TestRegister(t *testing.T) {
	UserID := "testuser"
	RequestJSON := `{"user_id":"testuser","email":"testuser@example.com","screen_name":"testuser1","password":"password"}`
	RegisteredResponse := fmt.Sprintf(`{"status_code": "201", "account_url": "/users/%s"}`, UserID)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/new", strings.NewReader(RequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h, err := SetUpDataBase()
	if err != nil {
		fmt.Println(err)
	}

	if assert.NoError(t, h.Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, RegisteredResponse, rec.Body.String())
	}
}

func TestLogin(t *testing.T) {
	UserID := "testuser"
	RequestJSON := `{"user_id":"testuser","password":"password"}`
	RegisteredResponse := fmt.Sprintf(`{"status_code": "201", "account_url": "/users/%s"}`, UserID)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/new", strings.NewReader(RequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h, err := SetUpDataBase()
	if err != nil {
		fmt.Println(err)
	}

	if assert.NoError(t, h.Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, RegisteredResponse, rec.Body.String())
	}
}
