package middleware

import (
	"context"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/AFK068/bot/internal/infrastructure/logger"

	handler "github.com/AFK068/bot/internal/infrastructure/handler/scrapper"
)

const (
	LinksPathPrefix = "/links"
	TgChatIDHeader  = "Tg-Chat-Id"
)

type UserChecker interface {
	CheckUserExistence(ctx context.Context, chatID int64) (bool, error)
}

func AuthLinkMiddleware(checker UserChecker, log *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if strings.HasPrefix(ctx.Path(), LinksPathPrefix) {
				log.Info("Request path starts with /links", "path", ctx.Path())

				if ok, err := checkAuthLink(ctx, checker, log); !ok {
					return err
				}
			}

			return next(ctx)
		}
	}
}

func checkAuthLink(ctx echo.Context, checker UserChecker, log *logger.Logger) (bool, error) {
	tgChatIDStr := ctx.Request().Header.Get(TgChatIDHeader)
	log.Info("Received Tg-Chat-Id header", "Tg-Chat-Id", tgChatIDStr)

	tgChatID, err := strconv.ParseInt(tgChatIDStr, 10, 64)
	if err != nil {
		log.Error("Failed to parse Tg-Chat-Id header", "Tg-Chat-Id", tgChatIDStr, "error", err)
		return false, handler.SendBadRequestResponse(ctx, handler.ErrInvalidRequestBody, handler.ErrDescriptionInvalidBody)
	}

	// Check that user exists in the repository.
	exist, err := checker.CheckUserExistence(ctx.Request().Context(), tgChatID)
	if err != nil {
		log.Error("Failed to check user existence", "Tg-Chat-Id", tgChatID, "error", err)
		return false, handler.SendBadRequestResponse(ctx, handler.ErrInternalError, handler.ErrDescriptionInternalError)
	}

	if !exist {
		log.Warn("User does not exist", "Tg-Chat-Id", tgChatID)
		return false, handler.SendUnauthorizedResponse(ctx, handler.ErrChatNotExist, handler.ErrDescriptionChatNotExist)
	}

	log.Info("User exists", "Tg-Chat-Id", tgChatID)

	return true, nil
}
