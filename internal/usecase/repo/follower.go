package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golanguzb70/udevslabs-twitter/config"
	"github.com/golanguzb70/udevslabs-twitter/internal/entity"
	"github.com/golanguzb70/udevslabs-twitter/pkg/logger"
	"github.com/golanguzb70/udevslabs-twitter/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
)

type FollowerRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewFollowerRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *FollowerRepo {
	return &FollowerRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *FollowerRepo) UpsertOrRemove(ctx context.Context, req entity.Follower) (entity.Follower, error) {

	query, args, err := r.pg.Builder.Insert("follower").
		Columns(`id, follower_id, following_id`).
		Values(uuid.NewString(), req.FollowerId, req.FollowingId).ToSql()
	if err != nil {
		return req, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {

		switch e := err.(type) {
		case *pgconn.PgError:
			// Handle PostgreSQL-specific errors
			switch e.Code {
			case "23505":
				query, args, err = r.pg.Builder.Delete("follower").Where(
					squirrel.Eq{
						"follower_id":  req.FollowerId,
						"following_id": req.FollowingId,
					}).ToSql()

				if err != nil {
					return req, err
				}

				_, err = r.pg.Pool.Exec(ctx, query, args...)
				if err == nil {
					req.UnFollowed = true
				}
			}
		}

		return entity.Follower{}, err
	}

	return req, nil
}

func (r *FollowerRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.UserList, error) {
	var (
		response             = entity.UserList{}
		createdAt, updatedAt time.Time
	)

	followingId := ""

	for i := 0; i < len(req.Filters); i++ {
		if req.Filters[i].Column == "following_id" {
			followingId = req.Filters[i].Value
			req.Filters = append(req.Filters[:i], req.Filters[i+1:]...)
			i--
		}
	}

	if followingId == "" {
		return response, fmt.Errorf("%sfollowing_id is required", "BAD_REQUEST")
	} else {
		req.Filters = append(req.Filters, entity.Filter{
			Column: "f.following_id",
			Type:   "eq",
			Value:  followingId,
		})
	}

	qeuryBuilder := r.pg.Builder.
		Select(`id, full_name, email, username, user_type, user_role, status, avatar_id, gender, created_at, updated_at`).
		From("follower f").Join("users as u ON u.id=f.follower_id")

	qeuryBuilder, where := PrepareGetListQuery(qeuryBuilder, req)

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, qeury, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.User
		err = rows.Scan(&item.ID, &item.FullName, &item.Email, &item.Username,
			&item.UserType, &item.UserRole, &item.Status, &item.AvatarId, &item.Gender, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("follower").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}
