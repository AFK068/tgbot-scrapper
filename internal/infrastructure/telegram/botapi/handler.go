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
		message := fmt.Sprintf("Link updated: %s", *linkUpdate.Url)
		if linkUpdate.Description != nil && *linkUpdate.Description != "" {
			message = fmt.Sprintf("%s\nDescription: %s", message, *linkUpdate.Description)
		}

		if linkUpdate.UserName != nil && *linkUpdate.UserName != "" {
			message = fmt.Sprintf("%s\nUpdated by: %s", message, *linkUpdate.UserName)
		}

		if linkUpdate.Type != nil {
			message = fmt.Sprintf("%s\nType: %s", message, *linkUpdate.Type)
		}

		if linkUpdate.СreatedAt != nil {
			message = fmt.Sprintf("%s\nCreated at: %s", message, linkUpdate.СreatedAt.Format("2006-01-02 15:04:05"))
		}

		h.Logger.Info("Sending message", "tgChatID", tgChatID, "message", message)
		h.Bot.SendMessage(tgChatID, message)
	}

	h.Logger.Info("Successfully processed PostUpdates request")

	return SendSuccessResponse(ctx, nil)
}
