package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
)

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

	h, err := SetUpDataBase()
	if err != nil {
		fmt.Println(err)
	}
	defer h.DB.Close()

	h.validate = validator.New()

	client := new(http.Client)
	server := httptest.NewServer(http.HandlerFunc(h.Register))
	defer server.Close()

	req, _ := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(RequestJSON))
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed: TestRegister(): %s", err.Error())
	}

	if res.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("Failed: TestRegister(): invaild content type. %s", err.Error())
	}

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("failed: Register() responsed different status code: %d", res.StatusCode)
	}
	if err := json.NewDecoder(res.Body).Decode(&Response); err != nil {
		t.Fatalf("failed: json.Unmarshal(): %s", err)
	}

	if Response.Status != "201" && Response.URL != "/users/testuser/" {
		t.Fatalf("failed: Register() responsed different body: %s", res.Body)
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

	h, err := SetUpDataBase()
	if err != nil {
		fmt.Println(err)
	}
	defer h.DB.Close()

	h.validate = validator.New()

	client := new(http.Client)
	server := httptest.NewServer(http.HandlerFunc(h.Login))
	defer server.Close()

	req, _ := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(RequestJSON))
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed TestLogin: %s", err.Error())
	}

	if res.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("Failed TestLogin: invaild content type, %s", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, res.Body)
		t.Fatalf("failed TestLogin: responsed different status code, %d", res.StatusCode)
	}
	if err := json.NewDecoder(res.Body).Decode(&Response); err != nil {
		t.Fatalf("failed: json.NewDecoder(): %s", err)
	}

	fmt.Println(Response)

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

	h, err := SetUpDataBase()
	if err != nil {
		fmt.Println(err)
	}
	defer h.DB.Close()

	h.validate = validator.New()

	t.Run("Not enough userid field", func(t *testing.T) {
		payload := new(ErrorPayload)

		RequestJSON := `{"userid":"","email":"testuser@example.com","screen_name":"testuser1","password":"password"}`

		client := new(http.Client)
		server := httptest.NewServer(http.HandlerFunc(h.Register))
		defer server.Close()

		req, _ := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(RequestJSON))
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed TestInvaildRegister/Not_enough_userid_field: %s", err.Error())
		}

		if err := json.NewDecoder(res.Body).Decode(payload); err != nil {
			t.Fatalf("Failed TestInvaildRegister/Not_enough_userid_field: %s", err.Error())
		}

		if res.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Failed TestInvaildRegister/Not_enough_userid_field: invaild content type, %s", res.Header.Get("Content-Type"))
		}

		if res.StatusCode != http.StatusBadRequest {
			fmt.Println(payload)
			t.Fatalf("Failed: invaild status code, %d", res.StatusCode)
		}

		if payload.StatusCode != "400" {
			t.Fatalf("Failed: invaild body")
		}

	})
	t.Run("Not enough email field", func(t *testing.T) {
		payload := new(ErrorPayload)
		RequestJSON := `{"userid":"","email":"testuser@example.com","screen_name":"testuser1","password":"password"}`

		client := new(http.Client)
		server := httptest.NewServer(http.HandlerFunc(h.Register))
		defer server.Close()

		req, _ := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(RequestJSON))
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed: %s", err.Error())
		}

		if err := json.NewDecoder(res.Body).Decode(payload); err != nil {
			t.Fatalf("Failed: %s", err.Error())
		}

		if res.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Failed: invaild content type, %s", res.Header.Get("Content-Type"))
		}

		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("failed: invaild status code, %d", res.StatusCode)
		}

		if payload.StatusCode != "400" {
			t.Fatalf("failed: invaild body")
		}
	})
	t.Run("Invaild email field", func(t *testing.T) {
		payload := new(ErrorPayload)
		RequestJSON := `{"userid":"testuser3","email":"testuser","screen_name":"testuser1","password":"password"}`

		client := new(http.Client)
		server := httptest.NewServer(http.HandlerFunc(h.Register))
		defer server.Close()

		req, _ := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(RequestJSON))
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed: %s", err.Error())
		}

		if err := json.NewDecoder(res.Body).Decode(payload); err != nil {
			t.Fatalf("Failed: %s", err.Error())
		}

		if res.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Failed: invaild content type, %s", res.Header.Get("Content-Type"))
		}
		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("failed: invaild status code, %d", res.StatusCode)
		}
		if payload.StatusCode != "400" {
			t.Fatalf("failed: invaild body, returned status_code is %s", payload.StatusCode)
		}
	})
	t.Run("Invaild userid field", func(t *testing.T) {
		payload := new(ErrorPayload)
		RequestJSON := `{"userid":"testtesttestuser3","email":"testuser@example.com","screen_name":"testuser1","password":"password"}`

		client := new(http.Client)
		server := httptest.NewServer(http.HandlerFunc(h.Register))
		defer server.Close()

		req, _ := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(RequestJSON))
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed: %s", err.Error())
		}

		if err := json.NewDecoder(res.Body).Decode(payload); err != nil {
			t.Fatalf("Failed: %s", err.Error())
		}

		if res.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Failed: invaild content type, %s", res.Header.Get("Content-Type"))
		}
		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("failed: Invaild userid field: invaild status code, %d", res.StatusCode)
		}
		if payload.StatusCode != "400" {
			t.Fatalf("failed: Invaild userid field: invaild body, status_code is %s", payload.StatusCode)
		}
	})
	t.Run("Invaild screen name field", func(t *testing.T) {
		payload := new(ErrorPayload)
		RequestJSON := `{"userid":"testtesttestuser3","email":"testuser@example.com","screen_name":"testtesttesttesttesttesttesttesttesttesttesttestuser1","password":"password"}`

		client := new(http.Client)
		server := httptest.NewServer(http.HandlerFunc(h.Register))
		defer server.Close()

		req, _ := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(RequestJSON))
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed: %s", err.Error())
		}

		if err := json.NewDecoder(res.Body).Decode(payload); err != nil {
			t.Fatalf("Failed: %s", err.Error())
		}

		if res.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Failed: invaild content type, %s", res.Header.Get("Content-Type"))
		}
		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("failed: Invaild screen name field: invaild status code, %d", res.StatusCode)
		}
		if payload.StatusCode != "400" {
			t.Fatalf("failed: Invaild screen name field: invaild body")
		}
	})
	t.Run("Invaild character in userid field", func(t *testing.T) {
		payload := new(ErrorPayload)
		RequestJSON := `{"userid":"test%&#*","email":"testuser@example.com","screen_name":"testuser1","password":"password"}`

		client := new(http.Client)
		server := httptest.NewServer(http.HandlerFunc(h.Register))
		defer server.Close()

		req, _ := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(RequestJSON))
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed: %s", err.Error())
		}

		if err := json.NewDecoder(res.Body).Decode(payload); err != nil {
			t.Fatalf("Failed: %s", err.Error())
		}

		if res.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Failed: invaild content type, %s", res.Header.Get("Content-Type"))
		}
		if res.StatusCode != http.StatusBadRequest {
			fmt.Println(payload)
			t.Fatalf("failed: Invaild character in userid field: invaild status code, %d", res.StatusCode)
		}
		if payload.StatusCode != "400" {
			t.Fatalf("failed: Invaild character in userid field: invaild body")
		}
	})
}
