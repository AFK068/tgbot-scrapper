package main

import (
	"github.com/AFK068/bot/config"
	"github.com/AFK068/bot/internal/application/scrapper"
	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/infrastructure/clients/bot"
	handler "github.com/AFK068/bot/internal/infrastructure/handler/scrapper"
	"github.com/AFK068/bot/internal/infrastructure/repository"
	"github.com/AFK068/bot/internal/infrastructure/server"
	"github.com/AFK068/bot/pkg/client/github"
	"github.com/AFK068/bot/pkg/client/stackoverflow"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			// Provide scrapper config.
			func() (*config.ScrapperConfig, error) {
				cfg, err := config.NewScrapperServerConfig(".")
				if err != nil {
					return nil, err
				}

				return cfg, nil
			},

			// Provide in-memory link repository.
			fx.Annotate(
				repository.NewInMemoryLinkRepository,
				fx.As(new(domain.ChatLinkRepository)),
			),

			// Provide scrapper handler.
			handler.NewScrapperHandler,

			// Provide stackoverflow client.
			stackoverflow.NewClient,

			// Provide github client.
			github.NewClient,

			// Provide bot client.
			func(cfg *config.ScrapperConfig) *bot.Client {
				return bot.NewClient(cfg.BotURL)
			},

			// Provide scrapper scheduler.
			scrapper.NewScrapperScheduler,

			// Provide scrapper server.
			server.NewScrapperServer,
		),
		fx.Invoke(
			// Run scrapper.
			func(s *server.ScrapperServer) error {
				return s.Start()
			},
		),
	).Run()
}
