package handler

import (
	"fmt"

	botapi "github.com/AFK068/bot/internal/api/openapi/bot/v1"
	"github.com/AFK068/bot/internal/application/bot"
	"github.com/labstack/echo/v4"
)

type BotHandler struct {
	Bot *bot.Bot
}

func NewBotHandler(b *bot.Bot) *BotHandler {
	return &BotHandler{
		Bot: b,
	}
}

func (h *BotHandler) PostUpdates(ctx echo.Context) error {
	var linkUpdate botapi.LinkUpdate
	if err := ctx.Bind(&linkUpdate); err != nil {
		return SendBadRequestResponse(ctx, ErrInvalidRequestBody, ErrDescriptionInvalidBody)
	}

	if linkUpdate.TgChatIds != nil {
		for _, tgChatID := range *linkUpdate.TgChatIds {
			if linkUpdate.Description != nil {
				h.Bot.SendMessage(tgChatID, *linkUpdate.Description)
			} else {
				h.Bot.SendMessage(tgChatID, fmt.Sprintf("Link updated: %s", *linkUpdate.Url))
			}
		}
	}

	return SendSuccessResponse(ctx, nil)
}
