package persona

import (
	"context"
	"log"

	authorization "github.com/eniehack/persona-server/gen/authorization"
)

// Authorization service example implementation.
// The example methods log the requests and return zero values.
type authorizationsrvc struct {
	logger *log.Logger
}

// NewAuthorization returns the Authorization service implementation.
func NewAuthorization(logger *log.Logger) authorization.Service {
	return &authorizationsrvc{logger}
}

// ログイン
func (s *authorizationsrvc) Login(ctx context.Context, p *authorization.LoginPayload) (err error) {
	s.logger.Print("authorization.login")
	return
}

// 新規登録
func (s *authorizationsrvc) Register(ctx context.Context, p *authorization.NewAccountPayload) (err error) {
	s.logger.Print("authorization.register")
	return
}
