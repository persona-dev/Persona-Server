package persona

import (
	"context"
	"log"

	post "github.com/eniehack/persona-server/gen/post"
)

// Post service example implementation.
// The example methods log the requests and return zero values.
type postsrvc struct {
	logger *log.Logger
}

// NewPost returns the Post service implementation.
func NewPost(logger *log.Logger) post.Service {
	return &postsrvc{logger}
}

// 新規投稿
func (s *postsrvc) Create(ctx context.Context, p *post.NewPostPayload) (err error) {
	s.logger.Print("post.create")
	return
}

// 投稿の参照
func (s *postsrvc) Reference(ctx context.Context, p *post.Post) (err error) {
	s.logger.Print("post.reference")
	return
}

// 投稿の削除
func (s *postsrvc) Delete(ctx context.Context, p *post.DeletePostPayload) (err error) {
	s.logger.Print("post.delete")
	return
}
