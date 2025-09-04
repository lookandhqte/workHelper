package inmemory

import (
	"sync"

	entity "github.com/lookandhqte/workHelper/internal/entity"
)

// MemoryStorage структура определяющая in-memory хранилище
type MemoryStorage struct {
	mu            sync.RWMutex
	account       *entity.Account
	token         *entity.Token
	lastAccountID int
}

// NewMemoryStorage создает новое хранилище in-memory
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		account:       &entity.Account{},
		token:         &entity.Token{},
		lastAccountID: 0,
	}
}
