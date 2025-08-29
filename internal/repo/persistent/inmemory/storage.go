package inmemory

import (
	"sync"

	entity "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

// MemoryStorage структура определяющая in-memory хранилище
type MemoryStorage struct {
	mu            sync.RWMutex
	account       *entity.Account
	lastAccountID int
}

// NewMemoryStorage создает новое хранилище in-memory
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		account:       &entity.Account{},
		lastAccountID: 0,
	}
}
