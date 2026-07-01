package transaction

import (
	"context"
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

func (u *UseCase) AddExpense(ctx context.Context, userID uuid.UUID,
	walletID uuid.UUID, updateID int64,
	amount domain.Money, category string) (*domain.Transaction, error) {

	if amount <= 0 {
		return nil, ErrAmountIsZero
	}

	if !domain.IsValidCategory(category) {
		return nil, ErrInvalidCategory
	}

	transactionId, err := uuid.NewV6()
	if err != nil {
		logger.Instance.Errorf("error while create transaction uuid: %v\n", err)
		return nil, ErrInternalWhileCreateTransaction
	}

	transaction := &domain.Transaction{
		ID:        transactionId,
		UserID:    userID,
		WalletID:  walletID,
		Amount:    amount,
		Category:  category,
		CreatedAt: time.Now(),
		UpdateID:  updateID,
	}

	err = u.repo.Save(ctx, transaction)
	if err != nil {
		if errors.Is(err, ErrTransactionAlreadyCreated) {
			return nil, err
		}

		logger.Instance.Errorf("error while save transaction: %v\n", err)
		return nil, ErrInternalWhileCreateTransaction
	}

	return transaction, nil
}
