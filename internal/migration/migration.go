package migration

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres" //nolint
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"

	"github.com/AFK068/bot/internal/config"
	"github.com/AFK068/bot/internal/infrastructure/logger"
)

func RunMigration(cfg *config.Config, log *logger.Logger) error {
	log.Info("Running migration")

	migrator, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.Migration.MigrationsPath),
		cfg.GetPostgresConnectionString(),
	)

	if err != nil {
		return fmt.Errorf("creating migrator: %w", err)
	}

	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("no migrations to apply")
			return nil
		}

		return fmt.Errorf("applying migrations: %w", err)
	}

	log.Info("migration completed")

	return nil
}
