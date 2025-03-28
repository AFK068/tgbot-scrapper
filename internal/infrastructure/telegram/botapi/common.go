package botapi

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"

	bottypes "github.com/AFK068/bot/internal/api/openapi/bot/v1"
)

const (
	ErrInvalidRequestBody = "invalid_request_body"
	ErrTgChatsIDIsEmpty   = "tg_chats_id_is_empty"
	ErrLinkIsEmpty        = "link_is_empty"

	ErrDescriptionInvalidBody      = "Invalid request body"
	ErrTgChatsIDIsEmptyDescription = "Tg chats id is empty"
	ErrLinkIsEmptyDescription      = "Link is empty"
)

func SendSuccessResponse(ctx echo.Context, data any) error {
	return ctx.JSON(http.StatusOK, data)
}

func SendBadRequestResponse(ctx echo.Context, err, description string) error {
	return ctx.JSON(http.StatusBadRequest, bottypes.ApiErrorResponse{
		Description:      aws.String(description),
		Code:             aws.String("400"),
		ExceptionMessage: aws.String(err),
	})
}
