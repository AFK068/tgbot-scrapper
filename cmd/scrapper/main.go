package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/AFK068/bot/internal/application/scrapper"
	"github.com/AFK068/bot/internal/config"
	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/infrastructure/clients/bot"
	"github.com/AFK068/bot/internal/infrastructure/httpapi/scrapperapi"
	"github.com/AFK068/bot/internal/infrastructure/logger"
	"github.com/AFK068/bot/internal/infrastructure/repository/link/ormrepo"
	"github.com/AFK068/bot/internal/infrastructure/repository/link/sqlrepo"
	"github.com/AFK068/bot/internal/infrastructure/server"
	"github.com/AFK068/bot/pkg/client/github"
	"github.com/AFK068/bot/pkg/client/stackoverflow"
	"github.com/AFK068/bot/pkg/txs"
)

const (
	ScrapperConfigPath = "config/scrapper.yaml"
	DBConfigPath       = "config/common.yaml"
)

func Postgres(dbConfig *config.Config, lc fx.Lifecycle) (domain.ChatLinkRepository, *pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(dbConfig.GetPostgresConnectionString())
	if err != nil {
		log.Fatal("failed to parse config: ", err)
	}

	ctx := context.Background()

	dbPool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		log.Fatal("failed to create database pool: ", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			dbPool.Close()
			return nil
		},
	})

	var repo domain.ChatLinkRepository
	if dbConfig.Storage.Type == domain.ORMRepository {
		repo = ormrepo.NewRepository(dbPool)
	} else {
		repo = sqlrepo.NewRepository(dbPool)
	}

	return repo, dbPool, nil
}

func main() {
	fx.New(
		fx.Provide(
			// Provide logger.
			logger.NewLogger,

			// Provide scrapper config.
			func() (*scrapper.Config, error) {
				cfg, err := scrapper.NewConfig(ScrapperConfigPath)
				if err != nil {
					return nil, err
				}

				return cfg, nil
			},

			// Provide db config.
			func() (*config.Config, error) {
				return config.NewConfig(DBConfigPath)
			},

			// Provide postgres repository.
			Postgres,

			// Provide transactor.
			fx.Annotate(
				txs.NewTxBeginner,
				fx.As(new(scrapperapi.Transactor)),
			),

			// Provide scrapper handler.
			scrapperapi.NewScrapperHandler,

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
			func(cfg *scrapper.Config, log *logger.Logger) bot.Service {
				return bot.NewClient(cfg.BotURL, log)
			},

			// Provide scrapper scheduler.
			scrapper.NewScrapperScheduler,

			// Provide scrapper server.
			server.NewScrapperServer,
		),
		fx.Invoke(
			// Run scrapper.
			func(s *server.ScrapperServer, lc fx.Lifecycle, log *logger.Logger) {
				s.RegisterHooks(lc, log)
			},
		),
	).Run()
}
