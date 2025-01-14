package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golanguzb70/udevslabs-twitter/config"
	"github.com/golanguzb70/udevslabs-twitter/internal/entity"
	"github.com/jackc/pgx/v4"
)

// Create/Delete Follower godoc
// @Router /follower [post]
// @Summary Create/Delete a follower
// @Description Create/Delete a follower
// @Security BearerAuth
// @Tags follower
// @Accept  json
// @Produce  json
// @Param follower body entity.Follower true "Follower object"
// @Success 200 {object} entity.Follower
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) FollowUnfollow(ctx *gin.Context) {
	var (
		body entity.Follower
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if ctx.GetHeader("user_type") == "user" {
		body.FollowerId = ctx.GetHeader("sub")
	}

	follower, err := h.UseCase.FollowerRepo.UpsertOrRemove(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating follower") {
		return
	}

	if follower.UnFollowed {
		userTags, err := h.UseCase.UserTagRepo.GetList(ctx, entity.GetListFilter{
			Page:  1,
			Limit: 1,
			Filters: []entity.Filter{
				{
					Column: "user_id",
					Type:   "eq",
					Value:  body.FollowerId,
				},
				{
					Column: "slug",
					Type:   "eq",
					Value:  body.FollowingId,
				},
			},
		})
		if h.HandleDbError(ctx, err, "error while getting user tag") {
			return
		}

		if len(userTags.Items) > 0 {
			err = h.UseCase.TagRepo.Delete(ctx, entity.Id{
				ID: userTags.Items[0].Id,
			})
			if h.HandleDbError(ctx, err, "error while getting user tag") {
				return
			}
		}
	} else {
		tag, err := h.UseCase.TagRepo.GetSingle(ctx, entity.Id{
			Slug: body.FollowingId,
		})
		if err == pgx.ErrNoRows {
			// create new tag
			tag, err = h.UseCase.TagRepo.Create(ctx, entity.Tag{
				Slug:  body.FollowingId,
				Level: 1,
			})
			if h.HandleDbError(ctx, err, "create tag") {
				return
			}
		}

		if tag.Id != "" {
			_, err = h.UseCase.UserTagRepo.Create(ctx, entity.UserTag{
				UserId: body.FollowerId,
				Tag:    tag,
			})
			if h.HandleDbError(ctx, err, "create user tag") {
				return
			}
		}
	}

	ctx.JSON(200, follower)
}

// GetFollowers godoc
// @Router /follower/list [get]
// @Summary Get a list of followers
// @Description Get a list of followers
// @Security BearerAuth
// @Tags follower
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param following_id query string false "following_id"
// @Param search query string false "search"
// @Success 200 {object} entity.UserList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetFollowers(ctx *gin.Context) {
	var (
		req entity.GetListFilter
	)

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("search", "")
	following_id := ctx.DefaultQuery("following_id", "")

	if ctx.GetHeader("user_type") == "user" {
		following_id = ctx.GetHeader("sub")
	}

	if following_id == "" {
		h.ReturnError(ctx, config.ErrorBadRequest, "following_id is required", http.StatusBadRequest)
		return
	}

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	if search != "" {
		req.Filters = append(req.Filters,
			entity.Filter{
				Column: "full_name",
				Type:   "search",
				Value:  search,
			},
			entity.Filter{
				Column: "username",
				Type:   "search",
				Value:  search,
			},
			entity.Filter{
				Column: "email",
				Type:   "search",
				Value:  search,
			},
		)
	}

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	users, err := h.UseCase.FollowerRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting users") {
		return
	}

	ctx.JSON(200, users)
}
