package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	OwnerWalletRole = "owner"
)

type Wallet struct {
	ID         uuid.UUID
	Title      string
	CreatedAt  time.Time
	InviteCode uuid.UUID
}
