package inmemory

import (
	"sync"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/pkg/cache"
)

type MemoryStorage struct {
	mu             sync.RWMutex
	accounts       map[int]*entity.Account
	integrations   map[int]*entity.Integration
	active_account *entity.Account
	lastAccountID  int
	cache          *cache.Cache
}

const (
	CACHE_EXPIRES_SEC = 604800
)

func NewMemoryStorage(c *cache.Cache) *MemoryStorage {
	return &MemoryStorage{
		accounts:       make(map[int]*entity.Account),
		integrations:   make(map[int]*entity.Integration),
		active_account: &entity.Account{},
		lastAccountID:  0,
		cache:          c,
	}
}
