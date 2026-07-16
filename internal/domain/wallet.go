package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	OwnerWalletRole  = "owner"
	MemberWalletRole = "member"
)

type Wallet struct {
	ID         uuid.UUID
	Title      string
	CreatedAt  time.Time
	InviteCode uuid.UUID
}

type WalletInfo struct {
	ID       uuid.UUID
	Title    string
	IsActive bool
}
