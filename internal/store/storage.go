package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetUserByID(context.Context, int64) (*User, error)
	}
	Comments interface {
		// get every comment from post with given post id
		GetByPostID(context.Context, int64) ([]Comment, error)
		// get specific comment by given comment id
		GetByComID(context.Context, int64) (*Comment, error)
		Create(context.Context, *Comment) error
		Delete(context.Context, int64) error
		Update(context.Context, *Comment) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
