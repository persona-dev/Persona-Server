// Code generated by goa v3.0.2, DO NOT EDIT.
//
// Authorization HTTP server types
//
// Command:
// $ goa gen github.com/eniehack/persona-server/design

package server

import (
	"unicode/utf8"

	authorization "github.com/eniehack/persona-server/gen/authorization"
	goa "goa.design/goa/v3/pkg"
)

// LoginRequestBody is the type of the "Authorization" service "login" endpoint
// HTTP request body.
type LoginRequestBody struct {
	// User ID
	Userid   *string `form:"userid,omitempty" json:"userid,omitempty" xml:"userid,omitempty"`
	Password *string `form:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
}

// RegisterRequestBody is the type of the "Authorization" service "register"
// endpoint HTTP request body.
type RegisterRequestBody struct {
	// User ID
	Userid *string `form:"userid,omitempty" json:"userid,omitempty" xml:"userid,omitempty"`
	// screen name of the user.
	ScreenName *string `form:"screen_name,omitempty" json:"screen_name,omitempty" xml:"screen_name,omitempty"`
	// password of the user.
	Password *string `form:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
}

// NewLoginPayload builds a Authorization service login endpoint payload.
func NewLoginPayload(body *LoginRequestBody) *authorization.LoginPayload {
	v := &authorization.LoginPayload{
		Userid:   *body.Userid,
		Password: *body.Password,
	}
	return v
}

// NewRegisterNewAccountPayload builds a Authorization service register
// endpoint payload.
func NewRegisterNewAccountPayload(body *RegisterRequestBody) *authorization.NewAccountPayload {
	v := &authorization.NewAccountPayload{
		Userid:     *body.Userid,
		ScreenName: *body.ScreenName,
		Password:   *body.Password,
	}
	return v
}

// ValidateLoginRequestBody runs the validations defined on LoginRequestBody
func ValidateLoginRequestBody(body *LoginRequestBody) (err error) {
	if body.Userid == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("userid", "body"))
	}
	if body.Password == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("password", "body"))
	}
	if body.Userid != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.userid", *body.Userid, "[^a-zA-Z0-9_]+"))
	}
	if body.Userid != nil {
		if utf8.RuneCountInString(*body.Userid) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.userid", *body.Userid, utf8.RuneCountInString(*body.Userid), 1, true))
		}
	}
	if body.Userid != nil {
		if utf8.RuneCountInString(*body.Userid) > 15 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.userid", *body.Userid, utf8.RuneCountInString(*body.Userid), 15, false))
		}
	}
	return
}

// ValidateRegisterRequestBody runs the validations defined on
// RegisterRequestBody
func ValidateRegisterRequestBody(body *RegisterRequestBody) (err error) {
	if body.Userid == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("userid", "body"))
	}
	if body.ScreenName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("screen_name", "body"))
	}
	if body.Password == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("password", "body"))
	}
	if body.Userid != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.userid", *body.Userid, "[^a-zA-Z0-9_]+"))
	}
	if body.Userid != nil {
		if utf8.RuneCountInString(*body.Userid) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.userid", *body.Userid, utf8.RuneCountInString(*body.Userid), 1, true))
		}
	}
	if body.Userid != nil {
		if utf8.RuneCountInString(*body.Userid) > 15 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.userid", *body.Userid, utf8.RuneCountInString(*body.Userid), 15, false))
		}
	}
	if body.ScreenName != nil {
		if utf8.RuneCountInString(*body.ScreenName) > 20 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.screen_name", *body.ScreenName, utf8.RuneCountInString(*body.ScreenName), 20, false))
		}
	}
	return
}