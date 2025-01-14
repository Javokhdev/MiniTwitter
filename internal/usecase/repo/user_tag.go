package repo

import (
	"context"
	"time"

	"github.com/golanguzb70/udevslabs-twitter/config"
	"github.com/golanguzb70/udevslabs-twitter/internal/entity"
	"github.com/golanguzb70/udevslabs-twitter/pkg/logger"
	"github.com/golanguzb70/udevslabs-twitter/pkg/postgres"
	"github.com/google/uuid"
)

type UserTagRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewUserTagRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UserTagRepo {
	return &UserTagRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *UserTagRepo) Create(ctx context.Context, req entity.UserTag) (entity.UserTag, error) {
	req.Id = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("user_tag").
		Columns(`id, user_id, tag_id`).
		Values(req.Id, req.UserId, req.Tag.Id).ToSql()
	if err != nil {
		return entity.UserTag{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.UserTag{}, err
	}

	return req, nil
}

func (r *UserTagRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.UserTagList, error) {
	var (
		response                   = entity.UserTagList{}
		createdAt, updatedAt       time.Time
		tagcreatedAt, tagupdatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`ut.id, ut.user_id, t.id, t.slug, t.level, t.created_at, t.updated_at, ut.created_at, ut.updated_at`).
		From("user_tag as ut").Join("tag as t ON ut.tag_id=t.id")

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
		var item entity.UserTag
		err = rows.Scan(&item.Id, &item.UserId, &item.Tag.Id, &item.Tag.Slug, &item.Tag.Level, &tagcreatedAt, &tagupdatedAt, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		item.Tag.CreatedAt = tagcreatedAt.Format(time.RFC3339)
		item.Tag.UpdatedAt = tagupdatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("user_tag").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *UserTagRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("user_tag").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserTagRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	qeury, args, err := r.pg.Builder.Update("user_tag").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
	if err != nil {
		return response, err
	}

	n, err := r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return response, err
	}

	response.RowsEffected = int(n.RowsAffected())

	return response, nil
}
