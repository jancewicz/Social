package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetUserByID(context.Context, int64) (*User, error)
		GetUserByEmail(context.Context, string) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string, inviteExp time.Duration) error
		Activate(ctx context.Context, token string) error
		Delete(ctx context.Context, userID int64) error
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
	Followers interface {
		Follow(ctx context.Context, followedID, userID int64) error
		Unfollow(ctx context.Context, followedID, userID int64) error
	}
	Roles interface {
		GetByName(ctx context.Context, role string) (*Role, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
		Roles:     &RolesStore{db},
	}
}

func withTransaction(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
