package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
	DBDSN    string
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	_ = godotenv.Load()
	err := cfg.parseEnv()
	return cfg, err
}

func (c *Config) parseEnv() error {
	botToken := os.Getenv("TUGRIK_BOT_TOKEN")
	if botToken == "" {
		return ErrorEmptyBotToken
	}

	dbDSN := os.Getenv("TUGRIK_DB_DSN")
	if dbDSN == "" {
		return ErrorEmptyDBDSN
	}

	c.BotToken = botToken
	c.DBDSN = dbDSN
	return nil
}
