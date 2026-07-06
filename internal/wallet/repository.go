package wallet

import (
	"context"

	"github.com/google/uuid"
	"github.com/whynullname/tugrikbot/internal/domain"
)

type Repository interface {
	GetUserWalletID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	CreateNewWallet(ctx context.Context, userID uuid.UUID, wallet *domain.Wallet) error
}
