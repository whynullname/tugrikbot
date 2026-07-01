package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/whynullname/tugrikbot/internal/domain"
	"github.com/whynullname/tugrikbot/internal/logger"
)

type UseCase struct {
	repo Repository
}

func NewUserUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) CreateUser(ctx context.Context, telegramUserId int64) error {
	isUserCreated, err := u.repo.IsUserCreated(ctx, telegramUserId)
	if err != nil {
		logger.Instance.Errorf("Error while get is user created: %v\n", err)
		return ErrInternalWhileCreateUser
	}

	if isUserCreated {
		return ErrUserAlreadyCreated
	}

	userID, err := uuid.NewV6()
	if err != nil {
		logger.Instance.Errorf("Error while create user uuid: %v\n", err)
		return ErrInternalWhileCreateUser
	}

	walletID, err := uuid.NewV6()
	if err != nil {
		logger.Instance.Errorf("Error while create wallet uuid: %v\n", err)
		return ErrInternalWhileCreateUser
	}

	currentTime := time.Now()
	user := &domain.User{
		ID:         userID,
		TelegramID: telegramUserId,
		CreatedAt:  currentTime,
	}

	walletTitle := fmt.Sprintf("personal_wallet_%d", telegramUserId)
	wallet := &domain.Wallet{
		ID:        walletID,
		Title:     walletTitle,
		CreatedAt: currentTime,
	}

	err = u.repo.CreateUser(ctx, user, wallet)
	if err != nil {
		logger.Instance.Errorf("Error while create user: %v\n", err)
		return ErrInternalWhileCreateUser
	}

	return nil
}

func (u *UseCase) GetUserID(ctx context.Context, telegramUserId int64) (uuid.UUID, error) {
	id, err := u.repo.GetUserID(ctx, telegramUserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, ErrUserNotFound
		}

		logger.Instance.Errorf("error while get user id: %v\n", err)
		return uuid.Nil, ErrInternalWhileGetUserID
	}

	return id, nil
}
