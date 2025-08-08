package inmemory

import (
	"fmt"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

//AddAccount создает аккаунт в in-memory хранилище
func (m *MemoryStorage) AddAccount(account *entity.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	account.ID = m.lastAccountID
	account.CacheExpires = account.CreatedAt + CacheExpires

	integrations := make([]entity.Integration, 0)
	m.accounts[account.ID].Integrations = integrations

	m.accounts[account.ID] = account
	m.lastAccountID++

	return nil
}

//GetAccounts возвращает все аккаунты из in-memory хранилища
func (m *MemoryStorage) GetAccounts() (*[]entity.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	accounts := make([]entity.Account, 0, len(m.accounts))
	for _, account := range m.accounts {
		accounts = append(accounts, *account)
	}

	return &accounts, nil
}

//GetAccount возвращает аккаунт по id из in-memory хранилища
func (m *MemoryStorage) GetAccount(id int) (*entity.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	account, exists := m.accounts[id]
	if !exists {
		return nil, fmt.Errorf("account not found to get account")
	}

	return account, nil
}

//GetAccountIntegrations возвращает интеграции аккаунта из in-memory хранилища
func (m *MemoryStorage) GetAccountIntegrations(accountID int) (*[]entity.Integration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	integrations := m.accounts[accountID].Integrations

	return &integrations, nil
}

//UpdateAccount обновляет аккаунт иsз in-memory хранилища
func (m *MemoryStorage) UpdateAccount(account *entity.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.accounts[account.ID]; !exists {
		return fmt.Errorf("account not found to update")
	}

	m.accounts[account.ID] = account
	return nil
}

//DeleteAccount удаляет аккаунт из in-memory хранилища
func (m *MemoryStorage) DeleteAccount(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.accounts[id]; !exists {
		return fmt.Errorf("account not found to delete")
	}

	delete(m.accounts, id)
	return nil
}

// func (m *MemoryStorage) SaveContacts(contact *[]entity.Contact) error {
// 	return nil
// }
