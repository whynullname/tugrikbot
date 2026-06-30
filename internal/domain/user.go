package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID
	TelegramID int64
	CreatedAt  time.Time
}
