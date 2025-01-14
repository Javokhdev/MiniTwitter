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
)

type AttachmentRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewAttachmentRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *AttachmentRepo {
	return &AttachmentRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *AttachmentRepo) Create(ctx context.Context, req entity.Attachment) (entity.Attachment, error) {
	req.Id = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("tweet_attachment").
		Columns(`id, tweet_id, filepath, content_type`).
		Values(req.Id, req.TweetId, req.FilePath, req.ContentType).ToSql()
	if err != nil {
		return entity.Attachment{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Attachment{}, err
	}

	return req, nil
}

func (r *AttachmentRepo) MultipleUpsert(ctx context.Context, req entity.AttachmentMultipleInsertRequest) ([]entity.Attachment, error) {
	hasNewAttachment := false

	tx, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	insertQuery := r.pg.Builder.Insert("tweet_attachment").
		Columns(`id, tweet_id, filepath, content_type`)

	for i, attachment := range req.Attachments {
		if attachment.Id == "" {
			hasNewAttachment = true

			attachment.Id = uuid.NewString()
			req.Attachments[i].Id = attachment.Id
			insertQuery = insertQuery.Values(attachment.Id, req.TweetId, attachment.FilePath, attachment.ContentType)
		}
	}

	existingAttachments := make(map[string]bool)
	for _, attachment := range req.Attachments {
		if attachment.Id != "" {
			existingAttachments[attachment.Id] = true
		}
	}

	if hasNewAttachment {
		query, args, err := insertQuery.ToSql()
		if err != nil {
			return nil, err
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			r.logger.Error("error while inserting tweet_attachment", err)
			return nil, err
		}
	}

	query, args, err := r.pg.Builder.Select("id").From("tweet_attachment").
		Where(squirrel.Eq{"tweet_id": req.TweetId}).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("error while getting ids of tweet_attachment", err)
		return nil, err
	}
	defer rows.Close()

	deletedAttachmentIds := []string{}

	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		if !existingAttachments[id] {
			deletedAttachmentIds = append(deletedAttachmentIds, id)
		}
	}

	for _, e := range deletedAttachmentIds {
		query, args, err := r.pg.Builder.Delete("tweet_attachment").Where("id = ?", e).ToSql()
		if err != nil {
			return nil, err
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			r.logger.Error("error while deleting tweet_attachment", err)
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Error("error while commiting tweet_attachment", err)
		return nil, err
	}

	attachments, err := r.GetList(ctx, entity.GetListFilter{
		Page:  1,
		Limit: 10,
		Filters: []entity.Filter{
			{
				Column: "tweet_id",
				Type:   "eq",
				Value:  req.TweetId,
			},
		},
	})
	if err != nil {
		r.logger.Error("error while getting tweet_attachment", err)
		return nil, err
	}

	return attachments.Items, nil
}

func (r *AttachmentRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Attachment, error) {
	response := entity.Attachment{}
	var (
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, tweet_id, filepath, content_type, created_at, updated_at`).
		From("tweet_attachment")

	switch {
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("id = ?", req.ID)
	default:
		return entity.Attachment{}, fmt.Errorf("GetSingle - invalid request")
	}

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.Attachment{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.Id, &response.TweetId, &response.FilePath, &response.ContentType, &createdAt, &updatedAt)
	if err != nil {
		return entity.Attachment{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *AttachmentRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.AttachmentList, error) {
	var (
		response             = entity.AttachmentList{}
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, tweet_id, filepath, content_type, created_at, updated_at`).
		From("tweet_attachment")

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
		var item entity.Attachment
		err = rows.Scan(&item.Id, &item.TweetId, &item.FilePath, &item.ContentType, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("tweet_attachment").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *AttachmentRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("tweet_attachment").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}
