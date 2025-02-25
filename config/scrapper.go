package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ScrapperConfig struct {
	Host   string `mapstructure:"SERVER_HOST"`
	Port   string `mapstructure:"SERVER_PORT"`
	BotURL string `mapstructure:"BOT_URL"`
}

func NewScrapperServerConfig(file string) (*ScrapperConfig, error) {
	viper.AddConfigPath(file)
	viper.SetConfigName("scrapper")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	config := &ScrapperConfig{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unmarshaling server config: %w", err)
	}

	return config, nil
}
