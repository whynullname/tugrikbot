package wallet

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/whynullname/tugrikbot/internal/domain"
	"github.com/whynullname/tugrikbot/internal/logger"
)

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) GetUserWalletID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	walletID, err := u.repo.GetUserWalletID(ctx, userID)
	if err != nil {
		logger.Instance.Errorf("error while get user wallet id: %v\n", err)
		return uuid.Nil, ErrInternalWhileGetUserWalletID
	}

	return walletID, nil
}

func (u *UseCase) CreateNewWallet(ctx context.Context, userID uuid.UUID, walletTitle string) (uuid.UUID, error) {
	walletID, err := uuid.NewV6()
	if err != nil {
		logger.Instance.Errorf("error while create uuid for new wallet: %v\n", err)
		return uuid.Nil, ErrInternalWhileCreateNewWallet
	}

	walletInviteCode, err := uuid.NewRandom()
	if err != nil {
		logger.Instance.Errorf("error while create uuid for wallet invite code: %v\n", err)
		return uuid.Nil, ErrInternalWhileCreateNewWallet
	}

	wallet := &domain.Wallet{
		ID:         walletID,
		Title:      walletTitle,
		CreatedAt:  time.Now(),
		InviteCode: walletInviteCode,
	}

	err = u.repo.CreateNewWallet(ctx, userID, wallet)
	if err != nil {
		logger.Instance.Errorf("error while create new wallet: %v\n", err)
		return uuid.Nil, ErrInternalWhileCreateNewWallet
	}

	return walletInviteCode, nil
}
