package wallet

import (
	"context"

	"github.com/google/uuid"
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
