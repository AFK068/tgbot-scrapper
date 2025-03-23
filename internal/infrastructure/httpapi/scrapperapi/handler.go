package scrapperapi

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"

	"github.com/AFK068/bot/internal/application/mapper"
	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/domain/apperrors"
	"github.com/AFK068/bot/internal/infrastructure/logger"
	"github.com/AFK068/bot/pkg/utils"

	scrappertypes "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
)

type ScrapperHandler struct {
	repository domain.ChatLinkRepository
	Logger     *logger.Logger
}

func NewScrapperHandler(repo domain.ChatLinkRepository, log *logger.Logger) *ScrapperHandler {
	return &ScrapperHandler{repository: repo, Logger: log}
}

// Register chat.
// (POST /tg-chat/{id}).
func (h *ScrapperHandler) PostTgChatId(ctx echo.Context, id int64) error { //nolint:revive,stylecheck // according to codgen interface
	h.Logger.Info("Registering chat", "ID", id)

	exist, err := h.repository.CheckUserExistence(ctx.Request().Context(), id)
	if err != nil {
		h.Logger.Error("Failed to check user existence", "ID", id, "error", err)
		return SendBadRequestResponse(ctx, ErrInternalError, ErrDescriptionInternalError)
	}

	if exist {
		h.Logger.Warn("Chat already exists", "ID", id)
		return SendBadRequestResponse(ctx, ErrChatAlreadyExist, ErrDescriptionChatAlreadyExist)
	}

	if err := h.repository.RegisterChat(ctx.Request().Context(), id); err != nil {
		h.Logger.Error("Failed to register chat", "ID", id, "error", err)
		return SendBadRequestResponse(ctx, ErrInternalError, ErrDescriptionInternalError)
	}

	h.Logger.Info("Successfully registered chat", "ID", id)

	return SendSuccessResponse(ctx, nil)
}

// Remove chat.
// (DELETE /tg-chat/{id}).
func (h *ScrapperHandler) DeleteTgChatId(ctx echo.Context, id int64) error { //nolint:revive,stylecheck // according to codgen interface
	h.Logger.Info("Removing chat", "ID", id)

	exist, err := h.repository.CheckUserExistence(ctx.Request().Context(), id)
	if err != nil {
		h.Logger.Error("Failed to check user existence", "ID", id, "error", err)
		return SendBadRequestResponse(ctx, ErrInternalError, ErrDescriptionInternalError)
	}

	if !exist {
		h.Logger.Warn("Chat does not exist", "ID", id)
		return SendNotFoundResponse(ctx, ErrChatNotExist, ErrDescriptionChatNotExist)
	}

	if err := h.repository.DeleteChat(ctx.Request().Context(), id); err != nil {
		h.Logger.Error("Failed to remove chat", "ID", id, "error", err)
		return SendBadRequestResponse(ctx, ErrInternalError, ErrDescriptionInternalError)
	}

	h.Logger.Info("Successfully removed chat", "ID", id)

	return SendSuccessResponse(ctx, nil)
}

// Add link tracking.
// (POST /links).
func (h *ScrapperHandler) PostLinks(ctx echo.Context, params scrappertypes.PostLinksParams) error {
	h.Logger.Info("Adding link for chat", "ID", params.TgChatId)

	var req scrappertypes.AddLinkRequest
	if err := ctx.Bind(&req); err != nil {
		h.Logger.Warn("Invalid request body", "error", err)
		return SendBadRequestResponse(ctx, ErrInvalidRequestBody, ErrDescriptionInvalidBody)
	}

	link, err := mapper.MapAddLinkRequestToDomain(params.TgChatId, &req)

	var linkValidateErr *apperrors.LinkValidateError
	if errors.As(err, &linkValidateErr) {
		h.Logger.Warn("Link validation error", "error", err)
		return SendBadRequestResponse(ctx, ErrInvalidRequestBody, ErrDescriptionInvalidBody)
	}

	var linkTypeErr *apperrors.LinkTypeError
	if errors.As(err, &linkTypeErr) {
		h.Logger.Warn("Link type not supported", "error", err)
		return SendBadRequestResponse(ctx, ErrLinkTypeNotSupported, ErrDescriptionLinkTypeNotSupported)
	}

	if err != nil {
		h.Logger.Error("Internal error", "error", err)
		return SendBadRequestResponse(ctx, ErrInternalError, ErrDescriptionInternalError)
	}

	if err := h.repository.SaveLink(ctx.Request().Context(), params.TgChatId, link); err != nil {
		h.Logger.Error("Failed to save link for chat", "ID", params.TgChatId, "error", err)
		return SendBadRequestResponse(ctx, ErrInternalError, ErrDescriptionInternalError)
	}

	h.Logger.Info("Successfully added link for chat", "ID", params.TgChatId)

	return SendSuccessResponse(ctx, nil)
}

// Remove link tracking.
// (DELETE /links).
func (h *ScrapperHandler) DeleteLinks(ctx echo.Context, params scrappertypes.DeleteLinksParams) error {
	h.Logger.Info("Removing link for chat", "ID", params.TgChatId)

	var req scrappertypes.RemoveLinkRequest
	if err := ctx.Bind(&req); err != nil {
		h.Logger.Warn("Invalid request body", "error", err)
		return SendBadRequestResponse(ctx, ErrInvalidRequestBody, ErrDescriptionInvalidBody)
	}

	if req.Link == nil || *req.Link == "" {
		h.Logger.Warn("Link is empty")
		return SendBadRequestResponse(ctx, ErrInvalidRequestBody, ErrDescriptionInvalidBody)
	}

	link := &domain.Link{
		URL: *req.Link,
	}

	err := h.repository.DeleteLink(ctx.Request().Context(), params.TgChatId, link)

	var linkNotExistErr *apperrors.LinkIsNotExistError
	if errors.As(err, &linkNotExistErr) {
		h.Logger.Warn("Link does not exist", "error", err)
		return SendNotFoundResponse(ctx, ErrLinkNotExist, ErrDescriptionLinkNotExist)
	}

	if err != nil {
		h.Logger.Error("Failed to remove link for chat", "ID", params.TgChatId, "error", err)
		return SendBadRequestResponse(ctx, ErrInternalError, ErrDescriptionInternalError)
	}

	h.Logger.Info("Successfully removed link for chat", "ID", params.TgChatId)

	return SendSuccessResponse(ctx, nil)
}

// Get all tracked links.
// (GET /links).
func (h *ScrapperHandler) GetLinks(ctx echo.Context, params scrappertypes.GetLinksParams) error {
	h.Logger.Info("Getting links for chat", "ID", params.TgChatId)

	links, err := h.repository.GetListLinks(ctx.Request().Context(), params.TgChatId)
	if err != nil {
		h.Logger.Error("Failed to get links for chat", "ID", params.TgChatId, "error", err)
		return SendBadRequestResponse(ctx, ErrInternalError, ErrDescriptionInternalError)
	}

	if len(links) == 0 {
		h.Logger.Info("No links found for chat", "ID", params.TgChatId)

		return SendSuccessResponse(ctx, scrappertypes.ListLinksResponse{
			Links: &[]scrappertypes.LinkResponse{},
			Size:  aws.Int32(0),
		})
	}

	linksResp := make([]scrappertypes.LinkResponse, len(links))
	for i, link := range links {
		linksResp[i] = scrappertypes.LinkResponse{
			Url:     aws.String(link.URL),
			Tags:    utils.SliceStringPtr(link.Tags),
			Filters: utils.SliceStringPtr(link.Filters),
		}
	}

	h.Logger.Info("Successfully retrieved links for chat", "ID", params.TgChatId)

	return SendSuccessResponse(ctx, scrappertypes.ListLinksResponse{
		Links: &linksResp,
		Size:  aws.Int32(int32(len(linksResp))), //nolint:gosec // as per the requirements
	})
}
