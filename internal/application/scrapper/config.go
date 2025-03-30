package scrapper

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Host   string `yaml:"host" env:"SCRAPPER_HOST" env-required:"true"`
	Port   string `yaml:"port" env:"SCRAPPER_PORT" env-required:"true"`
	BotURL string `yaml:"bot_url" env:"SCRAPPER_BOT_URL" env-required:"true"`
}

func NewConfig(file string) (*Config, error) {
	config := &Config{}

	if err := cleanenv.ReadConfig(file, config); err != nil {
		return nil, err
	}

	return config, nil
}
