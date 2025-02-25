package handler

import (
	"net/http"

	botapi "github.com/AFK068/bot/internal/api/openapi/bot/v1"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
)

const (
	ErrInvalidRequestBody     = "invalid_request_body"
	ErrDescriptionInvalidBody = "Invalid request body"
)

func SendSuccessResponse(ctx echo.Context, data any) error {
	return ctx.JSON(http.StatusOK, data)
}

func SendBadRequestResponse(ctx echo.Context, err, description string) error {
	return ctx.JSON(http.StatusBadRequest, botapi.ApiErrorResponse{
		Description:      aws.String(description),
		Code:             aws.String("400"),
		ExceptionMessage: aws.String(err),
	})
}
