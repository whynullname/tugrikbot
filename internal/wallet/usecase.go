package wallet

import (
	"context"
	"database/sql"
	"errors"
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

func (u *UseCase) GetActiveUserWalletID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	walletID, err := u.repo.GetActiveUserWalletID(ctx, userID)
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

func (u *UseCase) JoinToWallet(ctx context.Context, userID uuid.UUID, inviteCode string) error {
	walletInviteCode, err := uuid.Parse(inviteCode)
	if err != nil {
		return ErrInvalidInviteCode
	}

	err = u.repo.JoinToWallet(ctx, userID, walletInviteCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidInviteCode
		}

		if !errors.Is(err, ErrUserAlreadyInWallet) {
			logger.Instance.Errorf("error while join to wallet: %v\n", err)
			return ErrInternalWhileJoinToWallet
		}

		return err
	}

	return nil
}

func (u *UseCase) GetUserWallets(ctx context.Context, userID uuid.UUID) ([]domain.WalletInfo, error) {
	infos, err := u.repo.GetUserWallets(ctx, userID)
	if err != nil {
		logger.Instance.Errorf("error while get users wallets: %v\n", err)
		return nil, ErrInternalWhileGetWallets
	}

	return infos, nil
}

func (u *UseCase) SetActiveWallet(ctx context.Context, userID uuid.UUID, walletID string) error {
	walletUUID, err := uuid.Parse(walletID)
	if err != nil {
		return ErrInvalidWalletID
	}

	err = u.repo.SetActiveWallet(ctx, userID, walletUUID)
	if err != nil {
		if !errors.Is(err, ErrUserNotInWallet) {
			logger.Instance.Errorf("error while set active wallet: %v\n", err)
			return ErrInternalWhileSetActiveWallet
		}

		return err
	}

	return nil
}
