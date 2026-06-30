package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/whynullname/tugrikbot/internal/config"
	"github.com/whynullname/tugrikbot/internal/logger"
	"github.com/whynullname/tugrikbot/internal/postgres"
	"github.com/whynullname/tugrikbot/internal/telegram"
	"github.com/whynullname/tugrikbot/internal/transaction"
	"github.com/whynullname/tugrikbot/internal/user"
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

	pgDB, err := postgres.New(context.Background(), cfg.DBDSN)
	if err != nil {
		logger.Instance.Errorf("error connect to postgres: %v\n", err)
		return
	}
	defer pgDB.Close()

	transactionRepo := transaction.NewPostgresRepository(pgDB)
	transactionUseCase := transaction.NewUseCase(transactionRepo)

	userRepo := user.NewPostgresRepository(pgDB)
	userUseCase := user.NewUserUseCase(userRepo)

	middlewares := telegram.NewMiddlewares(userUseCase)
	bot, err := telegram.NewBot(cfg.BotToken, middlewares, transactionUseCase, userUseCase)
	if err != nil {
		logger.Instance.Errorf("error in initialize bot: %v\n", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	bot.Start(ctx)
}
