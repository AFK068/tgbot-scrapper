package server

import (
	"github.com/labstack/echo/v4"

	"github.com/AFK068/bot/config"
	"github.com/AFK068/bot/internal/application/bot"

	botapi "github.com/AFK068/bot/internal/api/openapi/bot/v1"
	handler "github.com/AFK068/bot/internal/infrastructure/handler/bot"
)

type BotServer struct {
	Config  *config.BotConfig
	Handler *handler.BotHandler
	Echo    *echo.Echo
	Bot     bot.Service
}

func NewBotServer(cfg *config.BotConfig, b bot.Service, hd *handler.BotHandler) *BotServer {
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
