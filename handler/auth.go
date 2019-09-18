package handler

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"encoding/json"
	"strconv"

	"golang.org/x/crypto/argon2"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/eniehack/persona-server/utils"
)

func MakeErrorResponseBody(statusCode int, detail string) []byte {
	json, err := json.Marshal(&ErrorPayload{
		StatusCode: strconv.Itoa(statusCode),
		Detail: detail,
	})
	if err != nil {
		log.Fatalf("MakeErrorResponseBody(): Failed to Marshal json: %s", err)
	}

	return json
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	requestData := new(LoginParams)
	
	if err := json.NewDecoder(r.Body).Decode(requestData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(MakeErrorResponseBody(http.StatusInternalServerError, "Internal Server Error. Please contact Admin."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.Validate.Struct(requestData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(MakeErrorResponseBody(http.StatusUnauthorized, "incorrect request json"))
		return
	}

	fmt.Println(requestData)

	if useridValidation := CheckRegexp(`[^a-zA-Z0-9_]+`, requestData.UserName); useridValidation {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(MakeErrorResponseBody(http.StatusUnauthorized, "invaild userid format"))
		return
	}

	UserID, Password, err := h.RoadPasswordAndUserID(requestData.UserName)

	match, err := comparePasswordAndHash(requestData.Password, Password)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(MakeErrorResponseBody(http.StatusInternalServerError, "Internal Server Error. Please contact Admin."))
		
		return
	}
	if !match {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(MakeErrorResponseBody(http.StatusUnauthorized, "incorrect userid or password"))
		fmt.Println(w)
		return 
	}

	if err := h.UpdateAt(UserID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(MakeErrorResponseBody(http.StatusInternalServerError, "Internal Server Error. Please contact Admin."))
		w.WriteHeader(http.StatusInternalServerError)
	}

	Token, err := GenerateJWTToken(UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(MakeErrorResponseBody(http.StatusInternalServerError, "Internal Server Error. Please contact Admin."))
		w.WriteHeader(http.StatusInternalServerError)
	}

	response, err := json.Marshal(&LoginResponseBody{
		Token: Token,
	})
	if err != nil {
		log.Fatalf("Login(): Failed json.Marshal() token: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write(MakeErrorResponseBody(http.StatusInternalServerError, "Internal Server Error. Please contact Admin."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

	return
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	requestData := new(RegisterParams)
	User := new(RegisterParams)
	w.Header().Set("Content-Type", "application/json")
	
	if err := json.NewDecoder(r.Body).Decode(requestData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(MakeErrorResponseBody(http.StatusInternalServerError, "Internal Server Error. Please contact Admin."))
		return
	}

	if err := h.Validate.Struct(requestData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(MakeErrorResponseBody(http.StatusBadRequest, "invaild request json format"))
		return
	}
	
	if useridValidation := CheckRegexp(`[^a-zA-Z0-9_]+`, requestData.UserID); useridValidation {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(MakeErrorResponseBody(http.StatusBadRequest, "invaild userid"))
		return
	}

	UserIDConflict, err := h.CheckUniqueUserID(strings.ToLower(requestData.UserID))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(MakeErrorResponseBody(http.StatusInternalServerError, "Internal Server Error. Please contact Admin."))
		return 
	}
	if !UserIDConflict {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(MakeErrorResponseBody(http.StatusInternalServerError, "Internal Server Error. Please contact Admin."))
		return
	}
	User.UserID = strings.ToLower(requestData.UserID)

	EMailConflict, err := h.CheckUniqueEmail(requestData.EMail)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(MakeErrorResponseBody(http.StatusInternalServerError, "Internal Server Error. Please contact Admin."))
		return
	}
	if !EMailConflict {
		w.WriteHeader(http.StatusConflict)
		w.Write(MakeErrorResponseBody(http.StatusConflict, "email address have already used"))
		return
	}
	User.EMail = requestData.EMail

	// 参考サイト(MIT License):https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go

	var p = &Argon2Params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	User.Password, err = generatePassword(requestData.Password, p)
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(MakeErrorResponseBody(http.StatusInternalServerError, "Internal Server Error. Please contact Admin."))
		return
	}

	if err := h.InsertUserData(User); err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(MakeErrorResponseBody(http.StatusInternalServerError, "Internal Server Error. Please contact Admin."))
		return
	}

	url := fmt.Sprintf("/users/%s/", User.UserID)

	returnjson, err := json.Marshal(map[string]string{
		"status_code": "201",
		"account_url": url,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(MakeErrorResponseBody(http.StatusInternalServerError, "Internal Server Error. Please contact Admin."))
		log.Fatalln(err)
		return 
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(returnjson)

	return
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
	var IsUnique sql.NullInt64
	db := h.DB

	BindParams := map[string]interface{}{
		"UserID": UserID,
	}

	Query, Params, err := sqlx.Named(
		`SELECT SUM(CASE WHEN user_id = :UserID THEN 1 ELSE 0 END) AS userid_count FROM users;`,
		BindParams,
	)
	if err != nil {
		return false, fmt.Errorf("Error CheckUniqueUserID(). Failed to set prepared statement: %s", err)
	}
	Rebind := db.Rebind(Query)

	if err := db.Get(&IsUnique, Rebind, Params...); err != nil {
		return false, fmt.Errorf("Error CheckUniqueUserID(). Failed to select user data: %s", err)
	}

	if IsUnique.Int64 == 0 {
		return true, nil
	}
	return false, nil
}

func (h *Handler) CheckUniqueEmail(EMail string) (bool, error) {
	var IsUnique sql.NullInt64
	db := h.DB
	BindParams := map[string]interface{}{
		"EMail": EMail,
	}
	Query, Params, err := sqlx.Named(
		`SELECT SUM(CASE WHEN user_id = :EMail THEN 1 ELSE 0 END) AS userid_count FROM users;`,
		BindParams,
	)
	if err != nil {
		return false, fmt.Errorf("Error CheckUniqueEMail(). Failed to set prepared statement: %s", err)
	}
	Rebind := db.Rebind(Query)

	if err := db.Get(&IsUnique, Rebind, Params...); err != nil {
		return false, fmt.Errorf("Error CheckUniqueEMail(). Failed to select user data: %s", err)
	}

	if IsUnique.Int64 == 0 {
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
		"Now":        time.Now().Format(time.RFC3339Nano),
		"Password":   User.Password,
	}
	Query, Params, err := sqlx.Named(
		"INSERT INTO users (user_id, email, screen_name, created_at, updated_at, password) VALUES (:UserID, :EMail, :ScreenName, :Now, :Now, :Password)",
		BindParams,
	)
	if err != nil {
		return fmt.Errorf("Error InsertUserData(). Failed to set prepared statement: %s", err)
	}

	Rebind := db.Rebind(Query)

	if _, err := db.Exec(Rebind, Params...); err != nil {
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
		"SELECT user_id, password FROM users WHERE user_id = :UserID OR email = :UserID",
		BindParams,
	)
	if err != nil {
		return "", "", fmt.Errorf("Error RoadPasswordAndUserID(). Failed to set prepared statement: %s", err)
	}

	Rebind := db.Rebind(Query)

	if db.QueryRowx(Rebind, Params...).Scan(&UserID, &Password); err != nil {
		return "", "", fmt.Errorf("Error RoadPasswordAndUserID(). Failed to select user data: %s", err)
	}
	return UserID, Password, nil
}

func (h *Handler) UpdateAt(RequestUserID string) error {
	db := h.DB
	BindParams := map[string]interface{}{
		"UserID": RequestUserID,
		"Now":    time.Now().Format(time.RFC3339Nano),
	}
	Query, Params, err := sqlx.Named(
		"UPDATE users SET updated_at = :Now WHERE user_id = :UserID",
		BindParams,
	)
	if err != nil {
		return fmt.Errorf("Error UpdateAt(). Failed to set prepared statement: %s", err)
	}

	Rebind := db.Rebind(Query)

	if _, err := db.Exec(Rebind, Params...); err != nil {
		return fmt.Errorf("Error UpdateAt(). Failed to update user data: %s", err)
	}
	return nil
}

func GenerateJWTToken(UserID string) (string, error) {
	PrivateKey, err := utils.LoadPrivateKey()
	if err != nil {
		return "", fmt.Errorf("LoadPrivateKey(): %s", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512,
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Add(time.Second * 5).Unix(),
			Audience:  UserID,
		},
	)
	t, err := token.SignedString(PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign string: %s", err)
	}
	return t, nil
}
