package wallet

import (
	"context"

	"github.com/google/uuid"
	"github.com/whynullname/tugrikbot/internal/domain"
)

type Repository interface {
	GetActiveUserWalletID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	CreateNewWallet(ctx context.Context, userID uuid.UUID, wallet *domain.Wallet) error
	JoinToWallet(ctx context.Context, userID uuid.UUID, walletInviteCode uuid.UUID) error
	GetUserWallets(ctx context.Context, userID uuid.UUID) ([]domain.WalletInfo, error)
	SetActiveWallet(ctx context.Context, userID uuid.UUID, walletID uuid.UUID) error
}
