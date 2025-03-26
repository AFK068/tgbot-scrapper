package main

import (
	"context"
	"os/signal"
	"syscall"

	"go.uber.org/fx"

	"github.com/AFK068/bot/internal/application/bot"
	"github.com/AFK068/bot/internal/infrastructure/clients/scrapper"
	"github.com/AFK068/bot/internal/infrastructure/logger"
	"github.com/AFK068/bot/internal/infrastructure/server"
	"github.com/AFK068/bot/internal/infrastructure/telegram/botapi"
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
			func() (*bot.Config, error) {
				botCfg, err := bot.NewConfig(ConfigPath)
				if err != nil {
					return nil, err
				}

				return botCfg, nil
			},

			// Provide scrapper client.
			func(cfg *bot.Config, log *logger.Logger) *scrapper.Client {
				return scrapper.NewClient(cfg.ScrapperURL, log)
			},

			// Provide bot.
			fx.Annotate(
				bot.NewBot,
				fx.As(new(bot.Service)),
			),

			// Provide bot handler.
			botapi.NewBotHandler,

			// Provide bot server.
			server.NewBotServer,
		),
		fx.Invoke(
			// Run bot.
			func(b bot.Service, log *logger.Logger) error {
				ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
				defer stop()

				err := b.Run(ctx)
				if err != nil {
					log.Error("Failed to run bot", "error", err)
				}

				return err
			},

			func(s *server.BotServer, lc fx.Lifecycle, log *logger.Logger) {
				s.RegisterHooks(lc, log)
			},
		),
	).Run()
}
