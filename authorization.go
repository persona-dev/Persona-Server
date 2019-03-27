package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/eniehack/simple-sns-go/app"
	"github.com/goadesign/goa"
)

// AuthorizationController implements the Authorization resource.
type AuthorizationController struct {
	*goa.Controller
}

// NewAuthorizationController creates a Authorization controller.
func NewAuthorizationController(service *goa.Service) *AuthorizationController {
	return &AuthorizationController{Controller: service.NewController("AuthorizationController")}
}

// Login runs the login action.
func (c *AuthorizationController) Login(ctx *app.LoginAuthorizationContext) error {
	// AuthorizationController_Login: start_implement

	// Put your logic here

	var userid string

	type Claims struct {
		Scope string `json:"scope"`
		jwt.StandardClaims
	}

	db, err := sql.Open("sqlite", "test.db")
	if err != nil {
		Log.Println(err)
		return ctx.InternalServerError()
	}

	defer db.Close()

	err := db.QueryRow(
		"SELECT user_id FROM users WHERE user_id = ? AND password = ?",
		ctx.Payload.Userid,
		ctx.Payload.Password,
	).Scan(&userid)

	if err != nil {
		log.Println(err)
		return ctx.BadRequest()
	}

	// TODO:JWTトークンの生成

	SigningKey := ""

	claims := Claims{
		"api:access",
		jwt.StandardClaims{
			Issuer: "simple-sns/Team-Ops",
			IssuedAt:  time.Now().Unix(),
			Audience:  userid,
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			NotBefore: time.Now().Add(-15 * time.Second).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(SigningKey)

	if err != nil {
		log.Println(err)
		return ctx.InternalServerError()
	}

	// TODO:ClaimにScopeを必ず付ける

	res := &app.Login{}
	res.Token := ss
	return ctx.OK(res)
	// AuthorizationController_Login: end_implement
}

// Register runs the register action.
func (c *AuthorizationController) Register(ctx *app.RegisterAuthorizationContext) error {
	// AuthorizationController_Register: start_implement

	// Put your logic here

	// TODO:DBにアクセス
	// TODO:ユーザーIDがすでに存在していないか確認(クエリを発行して実行に失敗した際のif err != nilでerr!=nilだった際のコードを書く？)

	return ctx.Created()
	// AuthorizationController_Register: end_implement
}
