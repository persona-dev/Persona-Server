package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
	"gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestMain(m *testing.M) {
	//h, err := SetUpDataBase()
	//db := h.DB
	ExitCode := m.Run()
	//db.Close()
	if err := os.Remove("test.db"); err != nil {
		fmt.Println(err)
	}
	os.Exit(ExitCode)
}

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
	type ResponseParams struct {
		URL    string `json:"account_url"`
		Status string `json:"status_code"`
	}
	Response := new(ResponseParams)
	RequestJSON := `{"userid":"testuser","email":"testuser@example.com","screen_name":"testuser1","password":"password"}`

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/new", strings.NewReader(RequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h, err := SetUpDataBase()
	if err != nil {
		fmt.Println(err)
	}
	defer h.DB.Close()

	if err := h.Register(c); err != nil {
		t.Fatalf("failed: Register(): %s", err)
	}

	if rec.Code != http.StatusCreated {
		t.Fatalf("failed: Register() responsed different status code: %d", rec.Code)
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &Response); err != nil {
		t.Fatalf("failed: json.Unmarshal(): %s", err)
	}

	if Response.Status != "201" && Response.URL != "/users/testuser/" {
		t.Fatalf("failed: Register() responsed different body: %s", rec.Body)
	}
}

func TestLogin(t *testing.T) {
	type ResponseParams struct {
		Token string `json:"token"`
	}
	type TokenPayload struct {
		ExpiresAt int64  `json:"exp"`
		IssuedAt  int64  `json:"iat"`
		NotBefore int64  `json:"nbf"`
		Audience  string `json:"aud"`
	}
	Response := new(ResponseParams)
	RequestJSON := `{"userid":"testuser","password":"password"}`

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(RequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("api/v1/auth/signature")

	h, err := SetUpDataBase()
	if err != nil {
		fmt.Println(err)
	}
	defer h.DB.Close()

	if err := h.Login(c); err != nil {
		t.Fatalf("failed: Login(): %s", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("failed: Login() responsed different status code: %d", rec.Code)
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &Response); err != nil {
		t.Fatalf("failed: json.Unmarshal(): %s", err)
	}

	Token := strings.Split(Response.Token, ".")
	Header := `{"alg":"RS512","typ":"JWT"}`
	ResponseHeader, err := base64.URLEncoding.DecodeString(Token[0])
	if err != nil {
		t.Fatalf("failed: base64.URLEncoding.DecodeString(): %s", err)
	}
	if string(ResponseHeader[:]) != Header {
		t.Fatalf("failed: Header is %s", ResponseHeader)
	}

	ResponsePayload, err := base64.URLEncoding.DecodeString(Token[1])
	if err != nil {
		t.Fatalf("failed: base64.URLEncoding.DecodeString(): %s", err)
	}

	Payload := new(TokenPayload)
	if err := json.Unmarshal(ResponsePayload, &Payload); err != nil {
		t.Fatalf("failed: json.Unmarshal(): %s", err)
	}

	if Payload.Audience != "testuser" {
		t.Fatalf("failed: Payload is %s", ResponsePayload)
	}
	if Payload.ExpiresAt < time.Now().Add(time.Minute*5).Unix() {
		t.Fatalf("failed: jwt token's exp invaild")
	}
}

func TestInvaildRegister(t *testing.T) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	h, err := SetUpDataBase()
	if err != nil {
		fmt.Println(err)
	}
	defer h.DB.Close()

	t.Run("Not enough userid field", func(t *testing.T) {
		RequestJSON := `{"userid":"","email":"testuser@example.com","screen_name":"testuser1","password":"password"}`
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/new", strings.NewReader(RequestJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		if err := h.Register(c); err != nil && err != echo.ErrBadRequest {
			t.Fatalf("failed: Register(): %s", err)
		}
	})
	t.Run("Not enough email field", func(t *testing.T) {
		RequestJSON := `{"userid":"testuser2","email":"","screen_name":"testuser1","password":"password"}`
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(RequestJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/new")
		if err := h.Register(c); err != nil && err != echo.ErrBadRequest {
			t.Fatalf("failed: Register(): %s", err)
		}
	})
	t.Run("Invaild email field", func(t *testing.T) {
		RequestJSON := `{"userid":"testuser3","email":"testuser","screen_name":"testuser1","password":"password"}`
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(RequestJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/new")
		if err := h.Register(c); err != nil && err != echo.ErrBadRequest {
			t.Fatalf("failed: Register(): %s", err)
		}
	})
	t.Run("Invaild userid field", func(t *testing.T) {
		RequestJSON := `{"userid":"testtesttestuser3","email":"testuser@example.com","screen_name":"testuser1","password":"password"}`
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(RequestJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/new")
		if err := h.Register(c); err != nil && err != echo.ErrBadRequest {
			t.Fatalf("failed: Register(): %s", err)
		}
	})
	t.Run("Invaild screen name field", func(t *testing.T) {
		RequestJSON := `{"userid":"testtesttestuser3","email":"testuser@example.com","screen_name":"testtesttesttesttesttesttesttesttesttesttesttestuser1","password":"password"}`
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(RequestJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/new")
		if err := h.Register(c); err != nil && err != echo.ErrBadRequest {
			t.Fatalf("failed: Register(): %s", err)
		}
	})
	t.Run("Invaild character in userid field", func(t *testing.T) {
		RequestJSON := `{"userid":"test%&#*/\","email":"testuser@example.com","screen_name":"testuser1","password":"password"}`
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(RequestJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/new")
		if err := h.Register(c); err != nil && err != echo.ErrBadRequest {
			t.Fatalf("failed: Register(): %s", err)
		}
	})
}
