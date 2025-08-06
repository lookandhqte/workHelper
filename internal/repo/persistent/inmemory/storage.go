package inmemory

import (
	"sync"

	entity "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	cache "git.amocrm.ru/gelzhuravleva/amocrm_golang/pkg/cache"
)

//MemoryStorage структура определяющая in-memory хранилище
type MemoryStorage struct {
	mu            sync.RWMutex
	accounts      map[int]*entity.Account
	integrations  map[int]*entity.Integration
	lastAccountID int
	cache         *cache.Cache
}

const (
	//CacheExpires константа, определяющая время жизни кэша
	CacheExpires = 604800
)

//NewMemoryStorage создает новое хранилище in-memory
func NewMemoryStorage(c *cache.Cache) *MemoryStorage {
	return &MemoryStorage{
		accounts:      make(map[int]*entity.Account),
		integrations:  make(map[int]*entity.Integration),
		lastAccountID: 0,
		cache:         c,
	}
}
