package middleware

import (
	"strconv"
	"strings"

	"github.com/AFK068/bot/internal/domain"
	handler "github.com/AFK068/bot/internal/infrastructure/handler/scrapper"
	"github.com/labstack/echo/v4"
)

const (
	LinksPathPrefix = "/links"
	TgChatIDHeader  = "Tg-Chat-Id"
)

func AuthLinkMiddleware(repo domain.ChatLinkRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if strings.HasPrefix(ctx.Path(), LinksPathPrefix) {
				if ok, err := checkAuthLink(ctx, repo); !ok {
					return err
				}
			}

			return next(ctx)
		}
	}
}

func checkAuthLink(ctx echo.Context, repo domain.ChatLinkRepository) (bool, error) {
	tgChatIDStr := ctx.Request().Header.Get(TgChatIDHeader)

	tgChatID, err := strconv.ParseInt(tgChatIDStr, 10, 64)
	if err != nil {
		return false, handler.SendBadRequestResponse(ctx, handler.ErrInvalidRequestBody, handler.ErrDescriptionInvalidBody)
	}

	// Check that user exists in the repository.
	if !repo.CheckUserExistence(tgChatID) {
		return false, handler.SendUnauthorizedResponse(ctx, handler.ErrChatNotExist, handler.ErrDescriptionChatNotExist)
	}

	return true, nil
}
