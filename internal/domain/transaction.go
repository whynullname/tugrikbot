package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Money int64

type Transaction struct {
	ID        uuid.UUID
	UserID    int64
	Amount    Money
	Category  string
	CreatedAt time.Time
}

var allowedCategory = []string{"еда", "транспорт", "развлечения", "жилье", "прочее"}

func IsValidCategory(category string) bool {
	for _, validCategory := range allowedCategory {
		if validCategory == category {
			return true
		}
	}

	return false
}

func (m Money) String() string {
	outputFirstPart := m / 100
	outputSecondPart := m % 100
	return fmt.Sprintf("%d.%02d ₽", outputFirstPart, outputSecondPart)
}
