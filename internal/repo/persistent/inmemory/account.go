package inmemory

import (
	"github.com/lookandhqte/workHelper/internal/entity"
)

// AddAccount создает аккаунт в in-memory хранилище
func (m *MemoryStorage) AddAccount(account *entity.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.account = account
	return nil
}

// GetAccount возвращает аккаунт по id из in-memory хранилища
func (m *MemoryStorage) GetAccount() (*entity.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.account, nil
}

// UpdateAccount обновляет аккаунт иsз in-memory хранилища
func (m *MemoryStorage) UpdateAccount(account *entity.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.account = account
	return nil
}

// DeleteAccount удаляет аккаунт из in-memory хранилища
func (m *MemoryStorage) DeleteAccount() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.account = &entity.Account{}
	m.lastAccountID--
	return nil
}
