package main

import (
	"github.com/AFK068/bot/config"
	"github.com/AFK068/bot/internal/application/bot"
	"github.com/AFK068/bot/internal/infrastructure/clients/scrapper"
	handler "github.com/AFK068/bot/internal/infrastructure/handler/bot"
	"github.com/AFK068/bot/internal/infrastructure/logger"
	"github.com/AFK068/bot/internal/infrastructure/server"
	"go.uber.org/fx"
)

const (
	ConfigPath = "."
)

func main() {
	fx.New(
		fx.Provide(
			// Provide logger.
			logger.NewLogger,

			// Provide bot config.
			func() (*config.BotConfig, error) {
				botCfg, err := config.NewBotConfig(ConfigPath)
				if err != nil {
					return nil, err
				}

				return botCfg, nil
			},

			// Provide scrapper client.
			func(cfg *config.BotConfig, log *logger.Logger) *scrapper.Client {
				return scrapper.NewClient(cfg.ScrapperURL, log)
			},

			// Provide bot.
			fx.Annotate(
				bot.NewBot,
				fx.As(new(bot.Service)),
			),

			// Provide bot handler.
			handler.NewBotHandler,

			// Provide bot server.
			server.NewBotServer,
		),
		fx.Invoke(
			// Run bot.
			func(b bot.Service) error {
				return b.Run()
			},

			// Start bot server.
			func(s *server.BotServer) error {
				return s.Start()
			},
		),
	).Run()
}
