package bot

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Token       string `yaml:"token" env:"BOT_TOKEN" env-required:"true"`
	Host        string `yaml:"host" env:"BOT_HOST" env-required:"true"`
	Port        string `yaml:"port" env:"BOT_PORT" env-required:"true"`
	ScrapperURL string `yaml:"scrapper_url" env:"BOT_SCRAPPER_URL" env-required:"true"`
}

func NewConfig(file string) (*Config, error) {
	config := &Config{}

	if err := cleanenv.ReadConfig(file, config); err != nil {
		return nil, err
	}

	return config, nil
}
