package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/whynullname/tugrikbot/internal/config"
	"github.com/whynullname/tugrikbot/internal/logger"
	"github.com/whynullname/tugrikbot/internal/telegram"
)

func main() {
	err := logger.InitializeLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Instance.Sync()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Instance.Errorf("error in read config: %v\n", err)
		return
	}

	bot, err := telegram.NewBot(cfg.BotToken)
	if err != nil {
		logger.Instance.Errorf("error in initialize bot: %v\n", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	bot.Start(ctx)
}
