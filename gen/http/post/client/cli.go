// Code generated by goa v3.0.2, DO NOT EDIT.
//
// Post HTTP client CLI support package
//
// Command:
// $ goa gen github.com/eniehack/persona-server/design

package client

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

	post "github.com/eniehack/persona-server/gen/post"
	goa "goa.design/goa/v3/pkg"
)

// BuildCreatePayload builds the payload for the Post create endpoint from CLI
// flags.
func BuildCreatePayload(postCreateBody string, postCreateToken string) (*post.NewPostPayload, error) {
	var err error
	var body CreateRequestBody
	{
		err = json.Unmarshal([]byte(postCreateBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, example of valid JSON:\n%s", "'{\n      \"body\": \"Et id qui.\"\n   }'")
		}
	}
	var token string
	{
		token = postCreateToken
	}
	v := &post.NewPostPayload{
		Body: body.Body,
	}
	v.Token = token
	return v, nil
}

// BuildReferencePayload builds the payload for the Post reference endpoint
// from CLI flags.
func BuildReferencePayload(postReferenceBody string, postReferencePostID string) (*post.Post, error) {
	var err error
	var body ReferenceRequestBody
	{
		err = json.Unmarshal([]byte(postReferenceBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, example of valid JSON:\n%s", "'{\n      \"body\": \"にゃーん\",\n      \"posted_at\": \"1993-10-29T23:49:50Z\",\n      \"screen_name\": \"ほげほげ\",\n      \"user_id\": \"hogehoge\"\n   }'")
		}
		err = goa.MergeErrors(err, goa.ValidateFormat("body.posted_at", body.PostedAt, goa.FormatDateTime))

		err = goa.MergeErrors(err, goa.ValidatePattern("body.user_id", body.UserID, "[^a-zA-Z0-9_]+"))
		if utf8.RuneCountInString(body.UserID) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.user_id", body.UserID, utf8.RuneCountInString(body.UserID), 1, true))
		}
		if utf8.RuneCountInString(body.UserID) > 15 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.user_id", body.UserID, utf8.RuneCountInString(body.UserID), 15, false))
		}
		if utf8.RuneCountInString(body.ScreenName) > 20 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.screen_name", body.ScreenName, utf8.RuneCountInString(body.ScreenName), 20, false))
		}
		if err != nil {
			return nil, err
		}
	}
	var postID string
	{
		postID = postReferencePostID
	}
	v := &post.Post{
		PostedAt:   body.PostedAt,
		UserID:     body.UserID,
		ScreenName: body.ScreenName,
		Body:       body.Body,
	}
	v.PostID = postID
	return v, nil
}

// BuildDeletePayload builds the payload for the Post delete endpoint from CLI
// flags.
func BuildDeletePayload(postDeletePostID string, postDeleteToken string) (*post.DeletePostPayload, error) {
	var postID string
	{
		postID = postDeletePostID
	}
	var token string
	{
		token = postDeleteToken
	}
	payload := &post.DeletePostPayload{
		PostID: postID,
		Token:  token,
	}
	return payload, nil
}