package transaction

import (
	"context"
	"sync"

	"github.com/whynullname/tugrikbot/internal/domain"
)

type InMemoryRepository struct {
	mutex        sync.Mutex
	transactions map[int64][]domain.Transaction
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{transactions: make(map[int64][]domain.Transaction)}
}

func (i *InMemoryRepository) Save(ctx context.Context, transaction *domain.Transaction) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.transactions[transaction.UserID] = append(i.transactions[transaction.UserID], *transaction)
	return nil
}
