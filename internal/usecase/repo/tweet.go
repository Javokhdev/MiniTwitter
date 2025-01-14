package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golanguzb70/udevslabs-twitter/config"
	"github.com/golanguzb70/udevslabs-twitter/internal/entity"
	"github.com/golanguzb70/udevslabs-twitter/pkg/logger"
	"github.com/golanguzb70/udevslabs-twitter/pkg/postgres"
	"github.com/google/uuid"
)

type TweetRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewTweetRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *TweetRepo {
	return &TweetRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *TweetRepo) Create(ctx context.Context, req entity.Tweet) (entity.Tweet, error) {
	req.Id = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("tweet").
		Columns(`id, owner_id, content, tags, status`).
		Values(req.Id, req.Owner.ID, req.Content, req.Tags, req.Status).ToSql()
	if err != nil {
		return entity.Tweet{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Tweet{}, err
	}

	return req, nil
}

func (r *TweetRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Tweet, error) {
	response := entity.Tweet{}
	var (
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, owner_id, content, tags, status, created_at, updated_at`).
		From("tweet")

	switch {
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("id = ?", req.ID)
	default:
		return entity.Tweet{}, fmt.Errorf("GetSingle - invalid request")
	}

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.Tweet{}, err
	}

	tags := []byte{}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.Id, &response.Owner.ID, &response.Content, &tags, &response.Status, &createdAt, &updatedAt)
	if err != nil {
		return entity.Tweet{}, err
	}

	err = json.Unmarshal(tags, &response.Tags)
	if err != nil {
		return entity.Tweet{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *TweetRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.TweetList, error) {
	var (
		response             = entity.TweetList{}
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, owner_id, content, status, created_at, updated_at, 
				(SELECT COALESCE(json_agg(row_to_json(ta)), '[]'::json) 
				 FROM tweet_attachment ta 
				 WHERE ta.tweet_id = tweet.id) AS attachments, 
				 (
					SELECT row_to_json(u)
					FROM users u
					WHERE u.id = tweet.owner_id
					LIMIT 1
				) AS user`).
		From("tweet")

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
		var item entity.Tweet
		var attachmentsJSON []byte
		var userJson []byte
		err = rows.Scan(&item.Id, &item.Owner.ID, &item.Content, &item.Status, &createdAt, &updatedAt, &attachmentsJSON, &userJson)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		err = json.Unmarshal(attachmentsJSON, &item.Attachments)
		if err != nil {
			return response, err
		}

		err = json.Unmarshal(userJson, &item.Owner)
		if err != nil {
			return response, err
		}

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("tweet").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *TweetRepo) Update(ctx context.Context, req entity.Tweet) (entity.Tweet, error) {
	mp := map[string]interface{}{
		"content":    req.Content,
		"status":     req.Status,
		"updated_at": "now()",
	}

	qeury, args, err := r.pg.Builder.Update("tweet").SetMap(mp).Where("id = ?", req.Id).ToSql()
	if err != nil {
		return entity.Tweet{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Tweet{}, err
	}

	return req, nil
}

func (r *TweetRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("tweet").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *TweetRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	qeury, args, err := r.pg.Builder.Update("tweet").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
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
