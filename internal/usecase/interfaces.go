// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/golanguzb70/udevslabs-twitter/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// UserRepo -.
	UserRepoI interface {
		Create(ctx context.Context, req entity.User) (entity.User, error)
		GetSingle(ctx context.Context, req entity.UserSingleRequest) (entity.User, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.UserList, error)
		Update(ctx context.Context, req entity.User) (entity.User, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}

	// SessionRepo -.
	SessionRepoI interface {
		Create(ctx context.Context, req entity.Session) (entity.Session, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Session, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.SessionList, error)
		Update(ctx context.Context, req entity.Session) (entity.Session, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}

	// Tag Repo
	TagRepoI interface {
		Create(ctx context.Context, req entity.Tag) (entity.Tag, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Tag, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.TagList, error)
		Update(ctx context.Context, req entity.Tag) (entity.Tag, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
		GetCategoriesFromDB(ctx context.Context) ([]string, error)
		TagTweetByContent(ctx context.Context, req entity.Tweet) (entity.Tweet, error)
	}

	// User Tag Repo
	UserTagRepoI interface {
		Create(ctx context.Context, req entity.UserTag) (entity.UserTag, error)
		Delete(ctx context.Context, req entity.Id) error
		GetList(ctx context.Context, req entity.GetListFilter) (entity.UserTagList, error)
	}

	// Follower Repo
	FollowerRepoI interface {
		UpsertOrRemove(ctx context.Context, req entity.Follower) (entity.Follower, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.UserList, error)
	}

	// Tweet attachment
	TweetAttachentRepoI interface {
		Create(ctx context.Context, req entity.Attachment) (entity.Attachment, error)
		MultipleUpsert(ctx context.Context, req entity.AttachmentMultipleInsertRequest) ([]entity.Attachment, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Attachment, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.AttachmentList, error)
		Delete(ctx context.Context, req entity.Id) error
	}

	// Tweet
	TweetI interface {
		Create(ctx context.Context, req entity.Tweet) (entity.Tweet, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Tweet, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.TweetList, error)
		Update(ctx context.Context, req entity.Tweet) (entity.Tweet, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}
)
