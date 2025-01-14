package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golanguzb70/udevslabs-twitter/config"
	"github.com/golanguzb70/udevslabs-twitter/internal/entity"
)

// CreateTag godoc
// @Router /tag [post]
// @Summary Create a tag
// @Description Create a tag
// @Security BearerAuth
// @Tags tag
// @Accept  json
// @Produce  json
// @Param tag body entity.Tag true "Tag object"
// @Success 200 {object} entity.Tag
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateTag(ctx *gin.Context) {
	var (
		body entity.Tag
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	tag, err := h.UseCase.TagRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating tag") {
		return
	}

	ctx.JSON(200, tag)
}

// GetTag godoc
// @Router /tag/{id} [get]
// @Summary Get a tag by ID
// @Description Get a tag by ID
// @Security BearerAuth
// @Tags tag
// @Accept  json
// @Produce  json
// @Param id path string true "Tag ID"
// @Success 200 {object} entity.Tag
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetTag(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	tag, err := h.UseCase.TagRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting tag") {
		return
	}

	ctx.JSON(200, tag)
}

// GetTags godoc
// @Router /tag/list [get]
// @Summary Get a list of users
// @Description Get a list of users
// @Security BearerAuth
// @Tags tag
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param search query string false "search"
// @Success 200 {object} entity.TagList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetTags(ctx *gin.Context) {
	var (
		req entity.GetListFilter
	)

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("search", "")

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	if search != "" {

		req.Filters = append(req.Filters,
			entity.Filter{
				Column: "slug",
				Type:   "search",
				Value:  search,
			},
		)
	}

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	tags, err := h.UseCase.TagRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting tag") {
		return
	}

	ctx.JSON(200, tags)
}

// UpdateTag godoc
// @Router /tag [put]
// @Summary Update a tag
// @Description Update a tag
// @Security BearerAuth
// @Tags tag
// @Accept  json
// @Produce  json
// @Param tag body entity.Tag true "Tag object"
// @Success 200 {object} entity.Tag
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateTag(ctx *gin.Context) {
	var (
		body entity.Tag
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	tag, err := h.UseCase.TagRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating tag") {
		return
	}

	ctx.JSON(200, tag)
}

// DeleteTag godoc
// @Router /tag/{id} [delete]
// @Summary Delete a tag
// @Description Delete a tag
// @Security BearerAuth
// @Tags tag
// @Accept  json
// @Produce  json
// @Param id path string true "Tag ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteTag(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	err := h.UseCase.TagRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting tag") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Tag deleted successfully",
	})
}
