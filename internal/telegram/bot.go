package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/whynullname/tugrikbot/internal/logger"
)

type TelegramBot struct {
	bot *bot.Bot
}

func NewBot(token string) (*TelegramBot, error) {
	b, err := bot.New(token, bot.WithDefaultHandler(handler))
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		bot: b,
	}, nil
}

func (t *TelegramBot) Start(ctx context.Context) {
	t.bot.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hello",
	})

	if err != nil {
		logger.Instance.Errorf("error in send message: %v\n", err)
	}
}
