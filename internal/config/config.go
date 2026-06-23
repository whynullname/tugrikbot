package config

import (
	"errors"
	"os"
)

type Config struct {
	BotToken string
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := cfg.parseEnv()
	return cfg, err
}

func (c *Config) parseEnv() error {
	if token := os.Getenv("TUGRIK_BOT_TOKEN"); token != "" {
		c.BotToken = token
		return nil
	}

	return errors.New("empty bot token in env")
}
