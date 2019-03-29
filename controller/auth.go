package controller

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var (
	ErrInvaildHash         = errors.New("the encoded hash is not in the correct format.")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type Claims struct {
	Scope string `json:"foo"`
	jwt.StandardClaims
}

type Argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func (h *Handler) Login(c echo.Context) error {

	userid := c.FormValue("userid")

	if userid == "" {
		return echo.ErrBadRequest
	}

	var password string
	db := h.DB
	defer db.Close()

	if err := db.QueryRow(
		"SELECT password FROM users WHERE user_id = ?",
		userid,
	).Scan(&password); err != nil {
		return echo.ErrBadRequest
	}
	// 送信されてきたpasswordの検証

	match, err := comparePasswordAndHash(c.FormValue("password"), password)
	if err != nil {
		log.Println(err)
	}
	if !match {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"status_code": "401",
		})
	}

	// updated_atの更新

	if err := db.QueryRow(
		"UPDATE users SET updated_at = ? WHERE user_id = ?",
		time.Now(),
		userid,
	); err != nil {
		return echo.ErrInternalServerError
	}

	// jwtの発行

	privateKey, err := ioutil.ReadFile("private-key.pem")
	if err != nil {
		log.Println("failed to road private key.", err)
		return echo.ErrInternalServerError
	}

	Key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		log.Println(err)
		return echo.ErrInternalServerError
	}

	claims := Claims{
		"api:access",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Add(15 * time.Second).Unix(),
			Audience:  userid,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	t, err := token.SignedString(Key)
	if err != nil {
		log.Println(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func (h *Handler) Register(c echo.Context) error {

	userid := c.FormValue("userid")

	//TODO:英数字のみであるか検証する
	if len := len(userid); CheckRegexp(`[^a-zA-Z0-9_]+`, userid) || len > 15 || len == 0 {
		return echo.ErrBadRequest
	}

	//TODO:英字はすべて小文字に変換する

	db := h.DB
	defer db.Close()

	if err := db.QueryRow(
		"SELECT user_id FROM users WHERE user_id = ?",
		strings.ToLower(userid),
	); err == nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"status_code": "409",
		})
	}

	// passwordをArgon2idで暗号化
	// 参考サイト(MIT License):https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go

	p := &Argon2Params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	password, err := generatePassword(c.FormValue("password"), p)

	if err != nil {
		return echo.ErrInternalServerError
	}
	// 指定されたデータをもとにINSERT

	if _, err := db.Exec(
		"INSERT INTO users (user_id, screen_name, created_at, updated_at, password) VALUES (?, ?, ?, ?, ?)",
		userid,
		c.FormValue("screen_name"),
		time.Now(),
		time.Now(),
		password,
	); err != nil {
		log.Println(err)
		return echo.ErrInternalServerError
	}

	url := fmt.Sprintf("/users/%s/", userid)

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
