package bot

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	ConfigBotName = "bot"
	ConfigBotType = "env"
)

type Config struct {
	Token       string `mapstructure:"BOT_TOKEN"`
	Host        string `mapstructure:"SERVER_HOST"`
	Port        string `mapstructure:"SERVER_PORT"`
	ScrapperURL string `mapstructure:"SCRAPPER_URL"`
}

func NewConfig(file string) (*Config, error) {
	viper.AddConfigPath(file)
	viper.SetConfigName(ConfigBotName)
	viper.SetConfigType(ConfigBotType)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unmarshaling server config: %w", err)
	}

	return config, nil
}
