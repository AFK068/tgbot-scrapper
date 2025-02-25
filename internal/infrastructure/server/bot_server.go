package server

import (
	"github.com/AFK068/bot/config"
	botapi "github.com/AFK068/bot/internal/api/openapi/bot/v1"
	"github.com/AFK068/bot/internal/bot"
	handler "github.com/AFK068/bot/internal/infrastructure/handler/bot"
	"github.com/labstack/echo/v4"
)

type BotServer struct {
	Config  *config.BotConfig
	Handler *handler.BotHandler
	Bot     *bot.Bot
	Echo    *echo.Echo
}

func NewBotServer(cfg *config.BotConfig, b *bot.Bot, hd *handler.BotHandler) *BotServer {
	return &BotServer{
		Config:  cfg,
		Handler: hd,
		Bot:     b,
		Echo:    echo.New(),
	}
}

func (s *BotServer) Start() error {
	botapi.RegisterHandlers(s.Echo, s.Handler)

	return s.Echo.Start(s.Config.Host + ":" + s.Config.Port)
}
