package usecase

import (
	"github.com/golanguzb70/udevslabs-twitter/config"
	"github.com/golanguzb70/udevslabs-twitter/internal/usecase/repo"
	"github.com/golanguzb70/udevslabs-twitter/pkg/logger"
	"github.com/golanguzb70/udevslabs-twitter/pkg/postgres"
)

// UseCase -.
type UseCase struct {
	UserRepo             UserRepoI
	SessionRepo          SessionRepoI
	TagRepo              TagRepoI
	UserTagRepo          UserTagRepoI
	FollowerRepo         FollowerRepoI
	TweetAttachmentsRepo TweetAttachentRepoI
	TweetRepo            TweetI
}

// New -.
func New(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		UserRepo:             repo.NewUserRepo(pg, config, logger),
		SessionRepo:          repo.NewSessionRepo(pg, config, logger),
		TagRepo:              repo.NewTagRepo(pg, config, logger),
		UserTagRepo:          repo.NewUserTagRepo(pg, config, logger),
		FollowerRepo:         repo.NewFollowerRepo(pg, config, logger),
		TweetAttachmentsRepo: repo.NewAttachmentRepo(pg, config, logger),
		TweetRepo:            repo.NewTweetRepo(pg, config, logger),
	}
}
