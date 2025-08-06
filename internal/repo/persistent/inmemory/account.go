package inmemory

import (
	"fmt"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

func (m *MemoryStorage) AddAccount(account *entity.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastAccountID++
	account.ID = m.lastAccountID
	account.CacheExpires = account.CreatedAt + CACHE_EXPIRES_SEC

	integrations := make([]entity.Integration, 0, 5)
	m.accounts[account.ID].Integrations = integrations

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

func (m *MemoryStorage) GetAccountIntegrations(accountID int) (*[]entity.Integration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	integrations := m.accounts[accountID].Integrations

	return &integrations, nil
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
