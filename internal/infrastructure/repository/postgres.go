package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/AFK068/bot/internal/config"
	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/infrastructure/repository/link/ormrepo"
	"github.com/AFK068/bot/internal/infrastructure/repository/link/sqlrepo"
)

func NewPostgresRepo(dbConfig *config.Config, lc fx.Lifecycle) (domain.ChatLinkRepository, *pgxpool.Pool, error) {
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
