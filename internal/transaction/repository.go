package transaction

import (
	"context"

	"github.com/whynullname/tugrikbot/internal/domain"
)

type Repository interface {
	Save(ctx context.Context, transaction *domain.Transaction) error
}
