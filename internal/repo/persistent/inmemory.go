package persistent

import (
	"fmt"
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
	BASE_ID_FOR_TOKENS  = 0
	ACCESS_EXPIRES_SEC  = 86400
	REFRESH_EXPIRES_SEC = 2592000
	CACHE_EXPIRES_SEC   = 604800
	BASE_URL            = "https://spetser.amocrm.ru/"
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

//Методы аккаунта

func (m *MemoryStorage) AddAccount(account *entity.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastAccountID++
	account.ID = m.lastAccountID
	account.CacheExpires = account.CreatedAt + CACHE_EXPIRES_SEC

	integrations := make([]entity.Integration, 0, 5)
	m.accounts[account.ID].Integrations = &integrations

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

func (m *MemoryStorage) ChangeActiveAccount(new_id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if new_id == m.active_account.ID {
		return fmt.Errorf("this acc is active")
	}
	account, err := m.GetAccount(new_id)
	if err != nil {
		return err
	}
	m.active_account = account
	return nil
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

//Методы интеграций

func (m *MemoryStorage) AddIntegration(integration *entity.Integration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.integrations[integration.AccountID] = integration
	integrationsPtr := m.accounts[integration.AccountID].Integrations

	*integrationsPtr = append(*integrationsPtr, *integration)
	return nil
}

func (m *MemoryStorage) GetIntegration(id int) (*entity.Integration, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	integration, exists := m.integrations[id]

	if !exists {
		return nil, fmt.Errorf("no integrations with these id")
	}

	return integration, nil

}

func (m *MemoryStorage) GetIntegrations() (*[]entity.Integration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	integrations := make([]entity.Integration, 0, len(m.integrations))
	for _, integration := range m.integrations {
		integrations = append(integrations, *integration)
	}

	return &integrations, nil
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

func (m *MemoryStorage) ReturnByClientID(client_id string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	integrations, _ := m.GetIntegrations()
	for id, integration := range *integrations {
		if integration.ClientID == client_id {
			return id, nil
		}
	}

	return 0, fmt.Errorf("haven't found here your integration")
}
