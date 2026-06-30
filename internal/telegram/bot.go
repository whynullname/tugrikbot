package telegram

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
	"github.com/whynullname/tugrikbot/internal/domain"
	"github.com/whynullname/tugrikbot/internal/logger"
	"github.com/whynullname/tugrikbot/internal/transaction"
	"github.com/whynullname/tugrikbot/internal/user"
)

type TelegramBot struct {
	bot                *bot.Bot
	middlewares        *Middlewares
	transactionUseCase *transaction.UseCase
	userUseCase        *user.UseCase
}

func NewBot(token string, middlewares *Middlewares, useCase *transaction.UseCase, userUseCase *user.UseCase) (*TelegramBot, error) {
	telegramBot := &TelegramBot{
		transactionUseCase: useCase,
		middlewares:        middlewares,
		userUseCase:        userUseCase,
	}

	opts := []bot.Option{
		bot.WithMiddlewares(middlewares.SaveUserId),
		bot.WithDefaultHandler(telegramBot.handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, err
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, telegramBot.startHandler)
	telegramBot.bot = b
	return telegramBot, nil
}

func (t *TelegramBot) Start(ctx context.Context) {
	t.bot.Start(ctx)
}

func (t *TelegramBot) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	err := t.userUseCase.CreateUser(ctx, update.Message.From.ID)
	if errors.Is(err, user.ErrUserAlreadyCreated) {
		t.sendMessage(ctx, b, update, "вы уже зарегистрировались")
		return
	}

	if errors.Is(err, user.ErrInternalWhileCreateUser) {
		t.sendMessage(ctx, b, update, "произошла внутреняя ошибка")
		return
	}

	t.sendMessage(ctx, b, update, "регистрация прошла успешно!")
}

func (t *TelegramBot) handler(ctx context.Context, b *bot.Bot, update *models.Update) {
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
	//TODO: вот тут доделать
	tr, err := t.transactionUseCase.AddExpense(ctx, uuid.Nil, uuid.Nil, money, messageTexts[1])
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
