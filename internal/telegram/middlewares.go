package telegram

import (
	"context"
	"errors"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/whynullname/tugrikbot/internal/domain"
	"github.com/whynullname/tugrikbot/internal/user"
)

type Middlewares struct {
	userUseCase *user.UseCase
}

const StartText = "/start"

func NewMiddlewares(userUseCase *user.UseCase) *Middlewares {
	return &Middlewares{userUseCase: userUseCase}
}

func (m *Middlewares) SaveUserId(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		if update.Message.From == nil {
			return
		}

		if update.Message.Text == StartText {
			next(ctx, bot, update)
			return
		}

		userId, err := m.userUseCase.GetUserID(ctx, update.Message.From.ID)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				SendMessage(ctx, bot, update, "сначала /start")
				return
			}

			SendMessage(ctx, bot, update, "произошла внутреняя ошибка")
			return
		}

		ctxWithUserId := context.WithValue(ctx, domain.UserIdContextKey, domain.UserContextID(userId))
		next(ctxWithUserId, bot, update)
	}
}
