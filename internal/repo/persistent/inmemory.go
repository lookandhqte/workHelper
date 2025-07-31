package persistent

import (
	"amocrm_golang/internal/entity"
	"amocrm_golang/pkg/cache"
	"fmt"
	"sync"
	"time"
)

type MemoryStorage struct {
	mu            sync.RWMutex
	accounts      map[int]*entity.Account
	integrations  map[int]*entity.Integration
	tokens        map[int]*entity.Token
	lastAccountID int
	cache         *cache.Cache
}

func NewMemoryStorage(c *cache.Cache) *MemoryStorage {
	return &MemoryStorage{
		accounts:      make(map[int]*entity.Account),
		integrations:  make(map[int]*entity.Integration),
		tokens:        make(map[int]*entity.Token),
		lastAccountID: 0,
		cache:         c,
	}
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
		return nil, fmt.Errorf("account not found")
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
		return fmt.Errorf("account not found")
	}

	m.accounts[account.ID] = account
	return nil
}

func (m *MemoryStorage) DeleteAccount(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.accounts[id]; !exists {
		return fmt.Errorf("account not found")
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
		return nil, fmt.Errorf("integration not found")
	}

	return integration, nil
}

func (m *MemoryStorage) UpdateIntegration(integration *entity.Integration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.integrations[integration.AccountID]; !exists {
		return fmt.Errorf("integration not found")
	}

	m.integrations[integration.AccountID] = integration
	return nil
}

func (m *MemoryStorage) DeleteIntegration(accountID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.integrations[accountID]; !exists {
		return fmt.Errorf("integration not found")
	}

	delete(m.integrations, accountID)
	return nil
}

//Функиця должна добавлять новые токены
func (m *MemoryStorage) AddTokens(response *entity.Token) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tokens[0] = response

	return nil
}

func (m *MemoryStorage) DeleteTokens() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.tokens, 0)

	return nil
}

//Функиця должна добавлять обновлять рефреш токен
func (m *MemoryStorage) UpdateRToken(refresh string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tokens[0].RefreshToken = refresh
	m.tokens[0].ExpiresIn = time.Now().Second() + 2592000 // 30 дней в секундах
	return nil
}

//Функиця должна добавлять обновлять access токен
func (m *MemoryStorage) UpdateAToken(access string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tokens[0].AccessToken = access
	m.tokens[0].ExpiresIn = time.Now().Second() + 86400 // 1 сутки в секундах

	return nil
}

func (m *MemoryStorage) GetRefreshToken() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.tokens) == 0 {
		return "", fmt.Errorf("no refresh key in storage")
	}
	return m.tokens[0].RefreshToken, nil
}
