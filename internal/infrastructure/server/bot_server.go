package server

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/AFK068/bot/internal/application/bot"
	"github.com/AFK068/bot/internal/infrastructure/logger"
	"github.com/AFK068/bot/internal/infrastructure/telegram/botapi"

	bottypes "github.com/AFK068/bot/internal/api/openapi/bot/v1"
)

type BotServer struct {
	Config  *bot.Config
	Handler *botapi.BotHandler
	Echo    *echo.Echo
	Bot     bot.Service
}

func NewBotServer(cfg *bot.Config, b bot.Service, hd *botapi.BotHandler) *BotServer {
	return &BotServer{
		Config:  cfg,
		Handler: hd,
		Bot:     b,
		Echo:    echo.New(),
	}
}

func (s *BotServer) Start() error {
	bottypes.RegisterHandlers(s.Echo, s.Handler)

	return s.Echo.Start(":" + s.Config.Port)
}

func (s *BotServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.Echo.Shutdown(ctx)
}

func (s *BotServer) RegisterHooks(lc fx.Lifecycle, log *logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Logger.Info("Starting bot server")

			go func() {
				if err := s.Start(); err != nil && err != http.ErrServerClosed {
					log.Error("Failed to start bot server", "error", err)
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			log.Logger.Info("Stopping bot server")

			if err := s.Stop(); err != nil {
				log.Error("Failed to stop bot server", "error", err)
			}

			return nil
		},
	})
}
