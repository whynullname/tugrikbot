package telegram

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/whynullname/tugrikbot/internal/domain"
	"github.com/whynullname/tugrikbot/internal/logger"
	"github.com/whynullname/tugrikbot/internal/transaction"
)

type TelegramBot struct {
	bot     *bot.Bot
	useCase *transaction.UseCase
}

func NewBot(token string, useCase *transaction.UseCase) (*TelegramBot, error) {
	telegramBot := &TelegramBot{useCase: useCase}
	b, err := bot.New(token, bot.WithDefaultHandler(telegramBot.handler))
	if err != nil {
		return nil, err
	}

	telegramBot.bot = b
	return telegramBot, nil
}

func (t *TelegramBot) Start(ctx context.Context) {
	t.bot.Start(ctx)
}

func (t *TelegramBot) handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.From == nil {
		return
	}

	messageTexts := strings.Fields(update.Message.Text)
	if len(messageTexts) < 2 {
		t.sendInvalidFormatMessage(ctx, b, update)
		return
	}

	parsedMoney, err := strconv.ParseInt(messageTexts[0], 10, 64)
	if err != nil {
		logger.Instance.Errorf("can't parse money: %v\n", err)
		t.sendInvalidFormatMessage(ctx, b, update)
		return
	}

	money := domain.Money(parsedMoney * 100)
	tr, err := t.useCase.AddExpense(ctx, update.Message.From.ID, money, messageTexts[1])
	if err != nil {
		if errors.Is(err, transaction.ErrInvalidCategory) {
			t.sendMessage(ctx, b, update, "неизвестная категория")
			return
		}

		if errors.Is(err, transaction.ErrAmountIsZero) {
			t.sendMessage(ctx, b, update, "потрачено должно быть больше 0")
			return
		}

		if errors.Is(err, transaction.ErrInternalWhileCreateTransaction) {
			t.sendMessage(ctx, b, update, "произошла системная ошибка при добавлении траты")
			return
		}

		logger.Instance.Errorf("internal error: %v\n", err)
		t.sendMessage(ctx, b, update, "произошла системная ошибка при добавлении траты")
		return
	}

	outputMessage := fmt.Sprintf("Записал %s на %s", tr.Amount, tr.Category)
	t.sendMessage(ctx, b, update, outputMessage)
}

func (t *TelegramBot) sendInvalidFormatMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	t.sendMessage(ctx, b, update, "не верный формат")
}

func (t *TelegramBot) sendMessage(ctx context.Context, b *bot.Bot, update *models.Update, message string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   message,
	})

	if err != nil {
		logger.Instance.Errorf("error in send message: %v\n", err)
	}
}
