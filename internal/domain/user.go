package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserContextKey string
type UserContextID uuid.UUID

const UserIdContextKey = UserContextKey("user_id")

type User struct {
	ID         uuid.UUID
	TelegramID int64
	CreatedAt  time.Time
}
