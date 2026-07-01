package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/whynullname/tugrikbot/internal/domain"
)

type Repository interface {
	IsUserCreated(ctx context.Context, telegramUserId int64) (bool, error)
	CreateUser(ctx context.Context, user *domain.User, wallet *domain.Wallet) error
	GetUserID(ctx context.Context, telegramUserID int64) (uuid.UUID, error)
}
