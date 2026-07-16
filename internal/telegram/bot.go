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
	"github.com/whynullname/tugrikbot/internal/user"
	"github.com/whynullname/tugrikbot/internal/wallet"
)

type TelegramBot struct {
	bot                *bot.Bot
	middlewares        *Middlewares
	transactionUseCase *transaction.UseCase
	userUseCase        *user.UseCase
	walletUseCase      *wallet.UseCase
}

func NewBot(token string, middlewares *Middlewares,
	transactionUseCase *transaction.UseCase, userUseCase *user.UseCase,
	walletUseCase *wallet.UseCase) (*TelegramBot, error) {

	telegramBot := &TelegramBot{
		transactionUseCase: transactionUseCase,
		middlewares:        middlewares,
		userUseCase:        userUseCase,
		walletUseCase:      walletUseCase,
	}

	opts := []bot.Option{
		bot.WithMiddlewares(middlewares.SaveUserId),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, err
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "expense", bot.MatchTypeCommandStartOnly, telegramBot.expenseHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "income", bot.MatchTypeCommandStartOnly, telegramBot.incomeHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommandStartOnly, telegramBot.startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "newwallet", bot.MatchTypeCommandStartOnly, telegramBot.newWalletHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "join", bot.MatchTypeCommandStartOnly, telegramBot.joinToWalletHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "wallets", bot.MatchTypeCommandStartOnly, telegramBot.getWallets)
	b.RegisterHandler(bot.HandlerTypeMessageText, "setwallet", bot.MatchTypeCommandStartOnly, telegramBot.setActiveWallet)
	b.RegisterHandler(bot.HandlerTypeMessageText, "balance", bot.MatchTypeCommandStartOnly, telegramBot.getBalanceHandler)
	telegramBot.bot = b
	return telegramBot, nil
}

func (t *TelegramBot) Start(ctx context.Context) {
	t.bot.Start(ctx)
}

func (t *TelegramBot) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	err := t.userUseCase.CreateUser(ctx, update.Message.From.ID)
	if errors.Is(err, user.ErrUserAlreadyCreated) {
		SendMessage(ctx, b, update, "вы уже зарегистрировались")
		return
	}

	if errors.Is(err, user.ErrInternalWhileCreateUser) {
		SendMessage(ctx, b, update, "произошла внутреняя ошибка")
		return
	}

	SendMessage(ctx, b, update, "регистрация прошла успешно!")
}

func (t *TelegramBot) newWalletHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := domain.GetUserIDByContext(ctx)
	messageTexts := strings.Fields(update.Message.Text)
	if len(messageTexts) < 2 {
		SendMessage(ctx, b, update, "введите имя кошелька")
		return
	}

	walletTitle := strings.Join(messageTexts[1:], " ")
	inviteCode, err := t.walletUseCase.CreateNewWallet(ctx, userID, walletTitle)
	if err != nil {
		if errors.Is(err, wallet.ErrInternalWhileCreateNewWallet) {
			SendMessage(ctx, b, update, "произошла системная ошибка при создании кошелька")
			return
		}

		SendMessage(ctx, b, update, "произошла внутреняя ошибка")
		return
	}

	outputMessage := fmt.Sprintf("Кошелек успешно создан, ключ для приглашения: %s", inviteCode)
	SendMessage(ctx, b, update, outputMessage)

}

func (t *TelegramBot) joinToWalletHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := domain.GetUserIDByContext(ctx)
	messageTexts := strings.Fields(update.Message.Text)
	if len(messageTexts) < 2 {
		SendMessage(ctx, b, update, "введите код приглашения")
		return
	}

	inviteCode := messageTexts[1]
	err := t.walletUseCase.JoinToWallet(ctx, userID, inviteCode)
	if err != nil {
		if errors.Is(err, wallet.ErrUserAlreadyInWallet) {
			SendMessage(ctx, b, update, "вы уже участник этого кошелька")
			return
		}

		if errors.Is(err, wallet.ErrInvalidInviteCode) {
			SendMessage(ctx, b, update, "не валидный код приглашения")
			return
		}

		SendMessage(ctx, b, update, "произошла внутреняя ошибка")
		return
	}

	SendMessage(ctx, b, update, "вы успешно вступили в кошелек")
}

func (t *TelegramBot) getWallets(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := domain.GetUserIDByContext(ctx)
	walletsInfos, err := t.walletUseCase.GetUserWallets(ctx, userID)
	if err != nil {
		SendMessage(ctx, b, update, "произошла внутреняя ошибка")
		return
	}

	walletsInfosMessage := make([]string, 0)
	for _, info := range walletsInfos {
		message := fmt.Sprintf("Название: %s, ID: %s, он выбран: %t", info.Title, info.ID, info.IsActive)
		walletsInfosMessage = append(walletsInfosMessage, message)
	}

	outputMessage := strings.Join(walletsInfosMessage, "\n")
	SendMessage(ctx, b, update, outputMessage)
}

func (t *TelegramBot) setActiveWallet(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := domain.GetUserIDByContext(ctx)
	messageTexts := strings.Fields(update.Message.Text)
	if len(messageTexts) < 2 {
		SendMessage(ctx, b, update, "введите id кошелька")
		return
	}

	walletId := messageTexts[1]
	err := t.walletUseCase.SetActiveWallet(ctx, userID, walletId)
	if err != nil {
		if errors.Is(err, wallet.ErrUserNotInWallet) {
			SendMessage(ctx, b, update, "вы не участвуете в этом кошельке")
			return
		}

		if errors.Is(err, wallet.ErrInvalidWalletID) {
			SendMessage(ctx, b, update, "не корректный id кошелька")
			return
		}

		SendMessage(ctx, b, update, "произошла внутреняя ошибка")
		return
	}

	SendMessage(ctx, b, update, "активный кошелек успешно изменен")
}

func (t *TelegramBot) getBalanceHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := domain.GetUserIDByContext(ctx)
	walletID, err := t.walletUseCase.GetActiveUserWalletID(ctx, userID)
	if err != nil {
		SendMessage(ctx, b, update, "произошла системная ошибка")
		return
	}

	money, err := t.transactionUseCase.GetWalletBalance(ctx, walletID)
	if err != nil {
		SendMessage(ctx, b, update, "произошла системная ошибка")
		return
	}

	outputMessage := fmt.Sprintf("Баланс: %s", money)
	SendMessage(ctx, b, update, outputMessage)
}

func (t *TelegramBot) expenseHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	t.transactionHandler(ctx, b, update, domain.Expense)
}

func (t *TelegramBot) incomeHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	t.transactionHandler(ctx, b, update, domain.Income)
}

func (t *TelegramBot) transactionHandler(ctx context.Context, b *bot.Bot, update *models.Update,
	transactionType domain.TransactionType) {
	messageTexts := strings.Fields(update.Message.Text)
	if len(messageTexts) < 3 {
		t.sendInvalidFormatMessage(ctx, b, update)
		return
	}

	parsedMoney, err := strconv.ParseInt(messageTexts[1], 10, 64)
	if err != nil {
		logger.Instance.Errorf("can't parse money: %v\n", err)
		t.sendInvalidFormatMessage(ctx, b, update)
		return
	}

	money := domain.Money(parsedMoney * 100)
	userID := domain.GetUserIDByContext(ctx)
	walletID, err := t.walletUseCase.GetActiveUserWalletID(ctx, userID)
	if err != nil {
		SendMessage(ctx, b, update, "произошла системная ошибка")
		return
	}

	tr, err := t.transactionUseCase.AddTransaction(ctx, userID, walletID, update.ID, money, messageTexts[2], transactionType)
	if err != nil {
		if errors.Is(err, transaction.ErrTransactionAlreadyCreated) {
			SendMessage(ctx, b, update, "транзакция уже добавлена")
			return
		}

		if errors.Is(err, transaction.ErrInvalidCategory) {
			SendMessage(ctx, b, update, "неизвестная категория")
			return
		}

		if errors.Is(err, transaction.ErrAmountIsZero) {
			SendMessage(ctx, b, update, "транзакция должна быть больше 0")
			return
		}

		logger.Instance.Errorf("internal error: %v\n", err)
		SendMessage(ctx, b, update, "произошла системная ошибка при добавлении транзакции")
		return
	}

	outputMessage := fmt.Sprintf("Транзакция %s успешно добавлена в %s", tr.Amount, tr.Category)
	SendMessage(ctx, b, update, outputMessage)
}

func (t *TelegramBot) sendInvalidFormatMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	SendMessage(ctx, b, update, "не верный формат")
}

func SendMessage(ctx context.Context, b *bot.Bot, update *models.Update, message string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   message,
	})

	if err != nil {
		logger.Instance.Errorf("error in send message: %v\n", err)
	}
}
