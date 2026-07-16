package domain

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Money int64
type TransactionType string

const (
	Expense TransactionType = "expense"
	Income  TransactionType = "income"
)

type Transaction struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	WalletID  uuid.UUID
	Amount    Money
	Category  string
	CreatedAt time.Time
	UpdateID  int64
	Type      TransactionType
}

var allowedExpenseCategory = []string{"еда", "транспорт", "развлечения", "жилье", "прочее"}
var allowedIncomeCategory = []string{"аванс", "зарплата", "прочее"}

func IsValidCategory(category string, transactionType TransactionType) bool {
	var categories []string

	if transactionType == Expense {
		categories = allowedExpenseCategory
	} else {
		categories = allowedIncomeCategory
	}

	return slices.Contains(categories, category)
}

func (m Money) String() string {
	outputFirstPart := m / 100
	outputSecondPart := m % 100
	return fmt.Sprintf("%d.%02d ₽", outputFirstPart, outputSecondPart)
}
