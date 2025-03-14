package main

import (
	"go.uber.org/fx"

	"github.com/AFK068/bot/config"
	"github.com/AFK068/bot/internal/application/scrapper"
	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/infrastructure/clients/bot"
	"github.com/AFK068/bot/internal/infrastructure/logger"
	"github.com/AFK068/bot/internal/infrastructure/repository"
	"github.com/AFK068/bot/internal/infrastructure/server"
	"github.com/AFK068/bot/pkg/client/github"
	"github.com/AFK068/bot/pkg/client/stackoverflow"

	handler "github.com/AFK068/bot/internal/infrastructure/handler/scrapper"
)

const (
	ConfigPath = "."
)

func main() {
	fx.New(
		fx.Provide(
			// Provide logger.
			logger.NewLogger,

			// Provide scrapper config.
			func() (*config.ScrapperConfig, error) {
				cfg, err := config.NewScrapperServerConfig(ConfigPath)
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
			fx.Annotate(
				stackoverflow.NewClient,
				fx.As(new(scrapper.StackOverlowQuestionFetcher)),
			),

			// Provide github client.
			fx.Annotate(
				github.NewClient,
				fx.As(new(scrapper.GitHubRepoFetcher)),
			),

			// Provide bot client.
			func(cfg *config.ScrapperConfig, log *logger.Logger) bot.Service {
				return bot.NewClient(cfg.BotURL, log)
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
