package handler

import (
	"fmt"

	botapi "github.com/AFK068/bot/internal/api/openapi/bot/v1"
	"github.com/AFK068/bot/internal/application/bot"
	"github.com/labstack/echo/v4"
)

type BotHandler struct {
	Bot bot.Service
}

func NewBotHandler(b bot.Service) *BotHandler {
	return &BotHandler{
		Bot: b,
	}
}

func (h *BotHandler) PostUpdates(ctx echo.Context) error {
	var linkUpdate botapi.LinkUpdate
	if err := ctx.Bind(&linkUpdate); err != nil {
		return SendBadRequestResponse(ctx, ErrInvalidRequestBody, ErrDescriptionInvalidBody)
	}

	if linkUpdate.TgChatIds == nil || len(*linkUpdate.TgChatIds) == 0 {
		return SendBadRequestResponse(ctx, ErrTgChatsIDIsEmpty, ErrTgChatsIDIsEmptyDescription)
	}

	if linkUpdate.Url == nil || *linkUpdate.Url == "" {
		return SendBadRequestResponse(ctx, ErrLinkIsEmpty, ErrLinkIsEmptyDescription)
	}

	for _, tgChatID := range *linkUpdate.TgChatIds {
		if linkUpdate.Description != nil && *linkUpdate.Description != "" {
			h.Bot.SendMessage(tgChatID, fmt.Sprintf("Link updated: %s\nDescription: %s", *linkUpdate.Url, *linkUpdate.Description))
		} else {
			h.Bot.SendMessage(tgChatID, fmt.Sprintf("Link updated: %s", *linkUpdate.Url))
		}
	}

	return SendSuccessResponse(ctx, nil)
}
