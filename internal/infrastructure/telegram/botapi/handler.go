package botapi

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/AFK068/bot/internal/application/bot"
	"github.com/AFK068/bot/internal/infrastructure/logger"

	bottypes "github.com/AFK068/bot/internal/api/openapi/bot/v1"
)

type BotHandler struct {
	Bot    bot.Service
	Logger *logger.Logger
}

func NewBotHandler(b bot.Service, l *logger.Logger) *BotHandler {
	return &BotHandler{
		Bot:    b,
		Logger: l,
	}
}

func (h *BotHandler) PostUpdates(ctx echo.Context) error {
	var linkUpdate bottypes.LinkUpdate
	if err := ctx.Bind(&linkUpdate); err != nil {
		h.Logger.Error("Failed to bind request body", "error", err)
		return SendBadRequestResponse(ctx, ErrInvalidRequestBody, ErrDescriptionInvalidBody)
	}

	if linkUpdate.TgChatIds == nil || len(*linkUpdate.TgChatIds) == 0 {
		h.Logger.Warn("TgChatIds is empty")
		return SendBadRequestResponse(ctx, ErrTgChatsIDIsEmpty, ErrTgChatsIDIsEmptyDescription)
	}

	if linkUpdate.Url == nil || *linkUpdate.Url == "" {
		h.Logger.Warn("Url is empty")
		return SendBadRequestResponse(ctx, ErrLinkIsEmpty, ErrLinkIsEmptyDescription)
	}

	for _, tgChatID := range *linkUpdate.TgChatIds {
		if linkUpdate.Description != nil && *linkUpdate.Description != "" {
			h.Logger.Info("Sending message with description", "tgChatID", tgChatID, "url", *linkUpdate.Url, "description", *linkUpdate.Description)
			h.Bot.SendMessage(tgChatID, fmt.Sprintf("Link updated: %s\nDescription: %s", *linkUpdate.Url, *linkUpdate.Description))
		} else {
			h.Logger.Info("Sending message without description", "tgChatID", tgChatID, "url", *linkUpdate.Url)
			h.Bot.SendMessage(tgChatID, fmt.Sprintf("Link updated: %s", *linkUpdate.Url))
		}
	}

	h.Logger.Info("Successfully processed PostUpdates request")

	return SendSuccessResponse(ctx, nil)
}
