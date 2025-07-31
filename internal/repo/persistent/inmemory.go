package persistent

import (
	"fmt"
	"sync"
	"time"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/pkg/cache"
)

type MemoryStorage struct {
	mu            sync.RWMutex
	accounts      map[int]*entity.Account
	integrations  map[int]*entity.Integration
	lastAccountID int
	cache         *cache.Cache
}

const (
	BASE_ID_FOR_TOKENS    = 0
	ACCESS_EXPIRES_SEC    = 86400
	REFRESH_EXPIRES_SEC   = 2592000
	CACHE_EXPIRES_SEC     = 604800
	REFRESH_THRESHOLD_SEC = 3600
)

func NewMemoryStorage(c *cache.Cache) *MemoryStorage {
	return &MemoryStorage{
		accounts:      make(map[int]*entity.Account),
		integrations:  make(map[int]*entity.Integration),
		lastAccountID: 0,
		cache:         c,
	}
}

func (m *MemoryStorage) GetConst(req string) (int, error) {
	switch req {
	case "id_tokens":
		return BASE_ID_FOR_TOKENS, nil
	case "access_exp":
		return ACCESS_EXPIRES_SEC, nil
	case "refresh_exp":
		return REFRESH_EXPIRES_SEC, nil
	case "cache_exp":
		return CACHE_EXPIRES_SEC, nil
	case "refresh_threshold":
		return REFRESH_THRESHOLD_SEC, nil
	}
	return 0, fmt.Errorf("no such constant")
}

func (m *MemoryStorage) AddAccount(account *entity.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastAccountID++
	account.ID = m.lastAccountID
	m.accounts[account.ID] = account

	return nil
}

func (m *MemoryStorage) GetAccounts() ([]*entity.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	accounts := make([]*entity.Account, 0, len(m.accounts))
	for _, account := range m.accounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (m *MemoryStorage) GetAccount(id int) (*entity.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	account, exists := m.accounts[id]
	if !exists {
		return nil, fmt.Errorf("account not found to get account")
	}

	return account, nil
}

func (m *MemoryStorage) GetAccountWithCache(id int) (*entity.Account, error) {
	if cached, ok := m.cache.Get(id); ok {
		return cached.(*entity.Account), nil
	}

	account, err := m.GetAccount(id)
	if err != nil {
		return nil, err
	}

	m.cache.Set(id, account, time.Hour)
	return account, nil
}

func (m *MemoryStorage) UpdateAccount(account *entity.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.accounts[account.ID]; !exists {
		return fmt.Errorf("account not found to update")
	}

	m.accounts[account.ID] = account
	return nil
}

func (m *MemoryStorage) DeleteAccount(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.accounts[id]; !exists {
		return fmt.Errorf("account not found to delete")
	}

	delete(m.accounts, id)
	return nil
}

func (m *MemoryStorage) AddIntegration(integration *entity.Integration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.integrations[integration.AccountID] = integration
	return nil
}

func (m *MemoryStorage) GetIntegrations() ([]*entity.Integration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	integrations := make([]*entity.Integration, 0, len(m.integrations))
	for _, integration := range m.integrations {
		integrations = append(integrations, integration)
	}

	return integrations, nil
}

func (m *MemoryStorage) GetAccountIntegrations(accountID int) (*entity.Integration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	integration, exists := m.integrations[accountID]
	if !exists {
		return nil, fmt.Errorf("integration not found to get account integrations")
	}

	return integration, nil
}

func (m *MemoryStorage) UpdateIntegration(integration *entity.Integration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.integrations[integration.AccountID]; !exists {
		return fmt.Errorf("integration not found to update")
	}

	m.integrations[integration.AccountID] = integration
	return nil
}

func (m *MemoryStorage) DeleteIntegration(accountID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.integrations[accountID]; !exists {
		return fmt.Errorf("integration not found to delete integration")
	}

	delete(m.integrations, accountID)
	return nil
}

func (m *MemoryStorage) AddTokens(response *entity.Token) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache.SetToken(BASE_ID_FOR_TOKENS, response, time.Duration(ACCESS_EXPIRES_SEC)*time.Second)

	return nil
}

func (m *MemoryStorage) DeleteTokens() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache.DeleteToken(BASE_ID_FOR_TOKENS)

	return nil
}

func (m *MemoryStorage) UpdateTokens(response *entity.Token) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache.Set(BASE_ID_FOR_TOKENS, response, time.Duration(ACCESS_EXPIRES_SEC)*time.Second)
	return nil
}

func (m *MemoryStorage) GetRefreshToken() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ref, exists := m.cache.GetToken(BASE_ID_FOR_TOKENS)
	if !exists {
		return "", fmt.Errorf("no refresh key in storage")
	}
	return ref.RefreshToken, nil
}

func (m *MemoryStorage) GetTokens() (*entity.Token, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Println(m)
	val, exists := m.cache.GetToken(BASE_ID_FOR_TOKENS)
	if !exists {
		return nil, fmt.Errorf("no tokens in storage")
	}
	return val, nil
}
