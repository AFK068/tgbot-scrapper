package handler

import (
	"errors"

	api "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/domain/apperrors"

	"github.com/AFK068/bot/pkg/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
)

// As for the 401 error, I asked the curator and he allowed it.

type ScrapperHandler struct {
	repository domain.ChatLinkRepository
}

func NewScrapperHandler(repo domain.ChatLinkRepository) *ScrapperHandler {
	return &ScrapperHandler{repository: repo}
}

// Register chat.
// (POST /tg-chat/{id}).
func (h *ScrapperHandler) PostTgChatId(ctx echo.Context, id int64) error { //nolint:revive,stylecheck // according to codgen interface
	if h.repository.CheckUserExistence(id) {
		return SendBadRequestResponse(ctx, ErrChatAlreadyExist, ErrDescriptionChatAlreadyExist)
	}

	if err := h.repository.RegisterChat(id); err != nil {
		return SendBadRequestResponse(ctx, ErrInternalError, ErrDescriptionInternalError)
	}

	return SendSuccessResponse(ctx, nil)
}

// Remove chat.
// (DELETE /tg-chat/{id}).
func (h *ScrapperHandler) DeleteTgChatId(ctx echo.Context, id int64) error { //nolint:revive,stylecheck // according to codgen interface
	if !h.repository.CheckUserExistence(id) {
		return SendNotFoundResponse(ctx, ErrChatNotExist, ErrDescriptionChatNotExist)
	}

	if err := h.repository.DeleteChat(id); err != nil {
		return SendBadRequestResponse(ctx, ErrInternalError, ErrDescriptionInternalError)
	}

	return SendSuccessResponse(ctx, nil)
}

// Add link tracking.
// (POST /links).
func (h *ScrapperHandler) PostLinks(ctx echo.Context, params api.PostLinksParams) error {
	var req api.AddLinkRequest
	if err := ctx.Bind(&req); err != nil {
		return SendBadRequestResponse(ctx, ErrInvalidRequestBody, ErrDescriptionInvalidBody)
	}

	if req.Link == nil || *req.Link == "" {
		return SendBadRequestResponse(ctx, ErrInvalidRequestBody, ErrDescriptionInvalidBody)
	}

	link := &domain.Link{
		URL: *req.Link,
	}

	if req.Tags != nil {
		link.Tags = *req.Tags
	}

	if req.Filters != nil {
		link.Filters = *req.Filters
	}

	if err := h.repository.SaveLink(params.TgChatId, link); err != nil {
		return SendBadRequestResponse(ctx, ErrInternalError, ErrDescriptionInternalError)
	}

	return SendSuccessResponse(ctx, nil)
}

// Remove link tracking.
// (DELETE /links).
func (h *ScrapperHandler) DeleteLinks(ctx echo.Context, params api.DeleteLinksParams) error {
	var req api.RemoveLinkRequest
	if err := ctx.Bind(&req); err != nil {
		return SendBadRequestResponse(ctx, ErrInvalidRequestBody, ErrDescriptionInvalidBody)
	}

	if req.Link == nil || *req.Link == "" {
		return SendBadRequestResponse(ctx, ErrInvalidRequestBody, ErrDescriptionInvalidBody)
	}

	link := &domain.Link{
		URL: *req.Link,
	}

	err := h.repository.DeleteLink(params.TgChatId, link)

	if err != nil {
		var linkNotExistErr *apperrors.LinkIsNotExistError
		if errors.As(err, &linkNotExistErr) {
			return SendNotFoundResponse(ctx, ErrLinkNotExist, ErrDescriptionLinkNotExist)
		}

		return SendBadRequestResponse(ctx, ErrInternalError, ErrDescriptionInternalError)
	}

	return SendSuccessResponse(ctx, nil)
}

// Get all tracked links.
// (GET /links).
func (h *ScrapperHandler) GetLinks(ctx echo.Context, params api.GetLinksParams) error {
	links, err := h.repository.GetListLinks(params.TgChatId)
	if err != nil {
		return SendBadRequestResponse(ctx, ErrInternalError, ErrDescriptionInternalError)
	}

	if len(links) == 0 {
		return SendSuccessResponse(ctx, api.ListLinksResponse{
			Links: &[]api.LinkResponse{},
			Size:  aws.Int32(0),
		})
	}

	linksResp := make([]api.LinkResponse, len(links))
	for i, link := range links {
		linksResp[i] = api.LinkResponse{
			Id:      aws.Int64(link.ID),
			Url:     aws.String(link.URL),
			Tags:    utils.SlicePtr(link.Tags),
			Filters: utils.SlicePtr(link.Filters),
		}
	}

	return SendSuccessResponse(ctx, api.ListLinksResponse{
		Links: &linksResp,
		Size:  aws.Int32(int32(len(linksResp))), //nolint:gosec // as per the requirements
	})
}
