package handler

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

func (h *Handler) Login(c echo.Context) error {

	if c.FormValue("userid") == "" {
		return echo.ErrBadRequest
	}

	UserID, Password, err := h.RoadPasswordAndUserID(c.FormValue("userid"))

	match, err := comparePasswordAndHash(c.FormValue("password"), Password)
	if err != nil {
		log.Println(err)
		return echo.ErrInternalServerError
	}
	if !match {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"status_code": "401",
		})
	}

	if err := h.UpdateAt(UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status_code": "500",
		})
	}

	Token, err := GenerateJWTToken(UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status_code": "500",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": Token,
	})
}

func (h *Handler) Register(c echo.Context) error {

	User := new(RegisterParams)
	User.ScreenName = c.FormValue("screen_name")

	if len := len(User.UserID); CheckRegexp(`[^a-zA-Z0-9_]+`, User.UserID) || len > 15 || len == 0 {
		return echo.ErrBadRequest
	}

	UserIDConflict, err := h.CheckUniqueUserID(strings.ToLower(c.FormValue("userid")))
	if err != nil {
		log.Println(err)
	}
	if !UserIDConflict {
		return c.JSON(http.StatusConflict, echo.Map{
			"status_code": "409",
		})
	}
	User.UserID = strings.ToLower(c.FormValue("userid"))

	EMailConflict, err := h.CheckUniqueEmail(c.FormValue("email"))
	if err != nil {
		log.Println(err)
	}
	if !EMailConflict {
		return c.JSON(http.StatusConflict, echo.Map{
			"status_code": "409",
		})
	}
	User.EMail = c.FormValue("email")

	// 参考サイト(MIT License):https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go

	var p = &Argon2Params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	User.Password, err = generatePassword(c.FormValue("password"), p)
	if err != nil {
		return echo.ErrInternalServerError
	}

	if err := h.InsertUserData(User); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status_code": "500",
		})
	}

	url := fmt.Sprintf("/users/%s/", User.UserID)

	return c.JSON(http.StatusCreated, echo.Map{
		"status_code": "201",
		"account_url": url,
	})
}

func generatePassword(password string, p *Argon2Params) (encodedHash string, err error) {

	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {

	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func comparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (p *Argon2Params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvaildHash
	}

	var version int
	if _, err := fmt.Sscanf(vals[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}

	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &Argon2Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}

func CheckRegexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}

func (h *Handler) CheckUniqueUserID(UserID string) (bool, error) {
	var IsUnique bool
	db := h.DB

	BindParams := map[string]interface{}{
		"UserID": UserID,
	}
	Query, Params, err := sqlx.Named(
		"SELECT user_id, CASE WHEN user_id=:UserID THEN 'false' ELSE 'true' FROM users;",
		BindParams,
	)
	if err != nil {
		return false, fmt.Errorf("Error CheckUniqueUserID(). Failed to set prepared statement: %s", err)
	}

	if err := db.QueryRowx(Query, Params).Scan(&IsUnique); err != nil {
		return false, fmt.Errorf("Error CheckUniqueUserID(). Failed to select user data: %s", err)
	}
	if IsUnique {
		return true, nil
	}
	return false, nil
}

func (h *Handler) CheckUniqueEmail(EMail string) (bool, error) {
	var IsUnique bool
	db := h.DB
	BindParams := map[string]interface{}{
		"EMail": EMail,
	}
	Query, Params, err := sqlx.Named(
		"SELECT email, CASE WHEN email=:EMail THEN 'false' ELSE 'true' FROM users;",
		BindParams,
	)
	if err != nil {
		return false, fmt.Errorf("Error CheckUniqueEMail(). Failed to set prepared statement: %s", err)
	}

	if err := db.QueryRowx(Query, Params).Scan(&IsUnique); err != nil {
		return false, fmt.Errorf("Error CheckUniqueEMail(). Failed to select user data: %s", err)
	}
	if IsUnique {
		return true, nil
	}
	return false, nil
}

func (h *Handler) InsertUserData(User *RegisterParams) error {
	db := h.DB
	BindParams := map[string]interface{}{
		"UserID":     User.UserID,
		"EMail":      User.EMail,
		"ScreenName": User.ScreenName,
		"Now":        time.Now(),
		"Password":   User.Password,
	}
	Query, Params, err := sqlx.Named(
		"INSERT INTO users (user_id, email, screen_name, created_at, updated_at, password) VALUES (:UserID, :EMail, :ScreenName, :Now, :Now, :Password)",
		BindParams,
	)
	if err != nil {
		return fmt.Errorf("Error InsertUserData(). Failed to set prepared statement: %s", err)
	}

	if _, err := db.Exec(Query, Params); err != nil {
		return fmt.Errorf("Error InsertUserData(). Failed to insert user data: %s", err)
	}
	return nil
}

func (h *Handler) RoadPasswordAndUserID(RequestUserID string) (string, string, error) {
	var UserID, Password string
	db := h.DB
	BindParams := map[string]interface{}{
		"UserID": RequestUserID,
	}
	Query, Params, err := sqlx.Named(
		"SELECT password, user_id FROM users WHERE user_id = :UserID OR email = :UserID",
		BindParams,
	)
	if err != nil {
		return "", "", fmt.Errorf("Error RoadPasswordAndUserID(). Failed to set prepared statement: %s", err)
	}

	if db.QueryRowx(Query, Params).Scan(&UserID, &Password); err != nil {
		return "", "", fmt.Errorf("Error RoadPasswordAndUserID(). Failed to select user data: %s", err)
	}
	return UserID, Password, nil
}

func (h *Handler) UpdateAt(RequestUserID string) error {
	db := h.DB
	BindParams := map[string]interface{}{
		"UserID": RequestUserID,
		"Now":    time.Now(),
	}
	Query, Params, err := sqlx.Named(
		"UPDATE users SET updated_at = :Now WHERE user_id = :UserID",
		BindParams,
	)
	if err != nil {
		return fmt.Errorf("Error UpdateAt(). Failed to set prepared statement: %s", err)
	}

	if _, err := db.Exec(Query, Params); err != nil {
		return fmt.Errorf("Error UpdateAt(). Failed to update user data: %s", err)
	}
	return nil
}

func LoadPrivateKey() ([]byte, error) {
	PrivateKey, err := ioutil.ReadFile("private-key.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to road private key: %s", err)
	}
	return PrivateKey, nil
}

func GenerateJWTToken(UserID string) (string, error) {
	PrivateKey, err := LoadPrivateKey()
	if err != nil {
		return "", fmt.Errorf("LoadPrivateKey(): %s", err)
	}

	Key, err := jwt.ParseRSAPrivateKeyFromPEM(PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse Privatekey: %s", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512,
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Add(time.Second * 5).Unix(),
			Audience:  UserID,
		},
	)
	t, err := token.SignedString(Key)
	if err != nil {
		return "", fmt.Errorf("failed to sign string: %s", err)
	}
	return t, nil
}
