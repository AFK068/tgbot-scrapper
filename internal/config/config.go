package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/AFK068/bot/internal/domain"
)

type Config struct {
	Storage   Storage   `yaml:"storage" env-required:"true"`
	Migration Migration `yaml:"migrations" env-required:"true"`
}

type Storage struct {
	Type         domain.RepositoryType `yaml:"type" env:"STORAGE_TYPE" env-required:"true" env-default:"sql"`
	Host         string                `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port         string                `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	DatabaseName string                `yaml:"database_name" env:"POSTGRES_DATABASE_NAME" env-required:"true"`
	User         string                `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password     string                `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
}

type Migration struct {
	MigrationsPath      string `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-required:"true"`
	MigrationsTableName string `yaml:"migrations_table_name" env:"MIGRATIONS_TABLE_NAME" env-default:"migrations"`
}

func NewConfig(filePath string) (*Config, error) {
	config := &Config{}

	if err := cleanenv.ReadConfig(filePath, config); err != nil {
		return nil, err
	}

	switch config.Storage.Type {
	case domain.DirectSQLRepository, domain.ORMRepository:
	default:
		config.Storage.Type = domain.DirectSQLRepository
	}

	return config, nil
}

func (cfg *Config) GetPostgresConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.DatabaseName,
	)
}
