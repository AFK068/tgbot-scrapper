package handler

import (
	"net/http"

	api "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
)

const (
	ErrChatAlreadyExist   = "chat_already_exists"
	ErrChatNotExist       = "chat_not_found"
	ErrInternalError      = "internal_error"
	ErrInvalidRequestBody = "invalid_request_body"

	ErrDescriptionChatAlreadyExist = "Chat already exists"
	ErrDescriptionChatNotExist     = "Chat not found"
	ErrDescriptionInternalError    = "Internal error"
	ErrDescriptionInvalidBody      = "Invalid request body"

	ErrLinkNotExist         = "link_not_exist"
	ErrLinkValidationError  = "link_validation_error"
	ErrLinkTypeNotSupported = "link_type_not_supported"

	ErrDescriptionLinkNotExist         = "Link not exist"
	ErrDescriptionLinkValidationError  = "Link validation error"
	ErrDescriptionLinkTypeNotSupported = "Link type not supported"
)

func SendSuccessResponse(ctx echo.Context, data any) error {
	return ctx.JSON(http.StatusOK, data)
}

func SendBadRequestResponse(ctx echo.Context, err, description string) error {
	return ctx.JSON(http.StatusBadRequest, api.ApiErrorResponse{
		Description:      aws.String(description),
		Code:             aws.String("400"),
		ExceptionMessage: aws.String(err),
	})
}

func SendNotFoundResponse(ctx echo.Context, err, description string) error {
	return ctx.JSON(http.StatusNotFound, api.ApiErrorResponse{
		Description:      aws.String(description),
		Code:             aws.String("404"),
		ExceptionMessage: aws.String(err),
	})
}

func SendUnauthorizedResponse(ctx echo.Context, err, description string) error {
	return ctx.JSON(http.StatusUnauthorized, api.ApiErrorResponse{
		Description:      aws.String(description),
		Code:             aws.String("401"),
		ExceptionMessage: aws.String(err),
	})
}
