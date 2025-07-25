package database

import (
	"amocrm_golang/cache"
	"amocrm_golang/model"
	"amocrm_golang/repository"
	"fmt"
	"sync"
	"time"
)

//Структура сдля работы с аккаунтами
type AccountService struct {
	repo repository.AccountRepository
}

//Структура для работы с интеграциями
type IntegrationService struct {
	repo repository.IntegrationRepository
}

//Создает новый экземпляр для работы с аккаунтами
func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

//Создает новый экземпляр для работы с интеграциями
func NewIntegrationService(repo repository.IntegrationRepository) *IntegrationService {
	return &IntegrationService{repo: repo}
}

//Создает новый аккаунт
func (s *AccountService) CreateAccount(account *model.Account) error {
	return s.repo.AddAccount(account)
}

//Возвращает все аккаунты
func (s *AccountService) GetAccountList() ([]*model.Account, error) {
	return s.repo.GetAccounts()
}

//Возвращает аккаунт по id
func (s *AccountService) GetAccountByID(id int) (*model.Account, error) {
	return s.repo.GetAccount(id)
}

//Обновляет существующий аккаунт
func (s *AccountService) UpdateAccount(account *model.Account) error {
	return s.repo.UpdateAccount(account)
}

//Удаляет аккаунт
func (s *AccountService) DeleteAccount(id int) error {
	return s.repo.DeleteAccount(id)
}

//Создает интеграцию
func (s *IntegrationService) CreateIntegration(integration *model.Integration) error {
	return s.repo.AddIntegration(integration)
}

//Возвращает все интеграции
func (s *IntegrationService) GetIntegrationList() ([]*model.Integration, error) {
	return s.repo.GetIntegrations()
}

//Возвращает все интеграции конкретного аккаунта
func (s *IntegrationService) GetAccountIntegrations(accountID int) (*model.Integration, error) {
	return s.repo.GetAccountIntegrations(accountID)
}

//Обновляет интеграцию
func (s *IntegrationService) UpdateIntegration(integration *model.Integration) error {
	return s.repo.UpdateIntegration(integration)
}

//Удаляет интеграцию
func (s *IntegrationService) DeleteIntegration(id int) error {
	return s.repo.DeleteIntegration(id)
}

//Структура in-memory хранилища
type MemoryStorage struct {
	mu            sync.RWMutex
	accounts      map[int]*model.Account
	integrations  map[int]*model.Integration
	lastAccountID int
	cache         *cache.Cache
}

// NewMemoryStorage создает новое in-memory хранилище.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		accounts:      make(map[int]*model.Account),
		integrations:  make(map[int]*model.Integration),
		lastAccountID: 0,
		cache:         cache.NewCache(),
	}
}

//Возвращает аккаунт по ID с использованием кэша
func (m *MemoryStorage) GetAccountWithCache(id int) (*model.Account, error) {
	if cached, ok := m.cache.Get(id); ok {
		return cached.(*model.Account), nil
	}

	account, err := m.GetAccount(id)
	if err != nil {
		return nil, err
	}

	m.cache.Set(id, account, time.Hour) // Кэшируем на 1 час
	return account, nil
}

//Обновляет аккаунт
func (m *MemoryStorage) UpdateAccount(account *model.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if account.ID == 0 {
		return fmt.Errorf("account ID is required")
	}

	if account.ID < 0 {
		return fmt.Errorf("invalid ID")
	}

	if _, exists := m.accounts[account.ID]; !exists {
		return fmt.Errorf("аккаунт не найден")
	}

	m.accounts[account.ID] = account
	return nil
}

//Удаляет аккаунт
func (m *MemoryStorage) DeleteAccount(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if id == 0 {
		return fmt.Errorf("ID is required")
	}
	if id < 0 {
		return fmt.Errorf("invalid ID")
	}

	if _, exists := m.accounts[id]; !exists {
		return fmt.Errorf("account not found")
	}

	delete(m.accounts, id)
	return nil
}

//Удаляет интеграцию по ID аккаунта
func (m *MemoryStorage) DeleteIntegration(accountID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if accountID == 0 {
		return fmt.Errorf("account ID is required")
	}

	if accountID < 0 {
		return fmt.Errorf("invalid ID")
	}

	if _, exists := m.integrations[accountID]; !exists {
		return fmt.Errorf("integration not found")
	}

	delete(m.integrations, accountID)
	return nil
}

//Обновляет интеграцию
func (m *MemoryStorage) UpdateIntegration(integration *model.Integration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.integrations[integration.AccountID]; !exists {
		return fmt.Errorf("интеграция не найдена")
	}

	m.integrations[integration.AccountID] = integration
	return nil
}

//Добавляет аккаунт
func (m *MemoryStorage) AddAccount(account *model.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastAccountID++
	account.ID = m.lastAccountID
	m.accounts[account.ID] = account

	return nil
}

//Возвращает все аккаунты
func (m *MemoryStorage) GetAccounts() ([]*model.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	accounts := make([]*model.Account, 0, len(m.accounts))
	for _, account := range m.accounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

//Возвращает аккаунт по ID
func (m *MemoryStorage) GetAccount(id int) (*model.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var account model.Account
	if _, exists := m.accounts[id]; !exists {
		return nil, fmt.Errorf("аккаунт не найден по id")
	}
	account = *m.accounts[id]
	return &account, nil
}

//Возвращает интеграции аккаунта
func (m *MemoryStorage) GetAccountIntegrations(accountID int) (*model.Integration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if accountID == 0 {
		return nil, fmt.Errorf("account ID is required")
	}
	if accountID < 0 {
		return nil, fmt.Errorf("invalid ID")
	}

	integration, exists := m.integrations[accountID]

	if !exists {
		return nil, fmt.Errorf("integration for account %d not found", accountID)
	}
	return integration, nil
}

//Добавляет интеграцию
func (m *MemoryStorage) AddIntegration(integration *model.Integration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.integrations[integration.AccountID] = integration
	return nil
}

//Возвращает интеграции
func (m *MemoryStorage) GetIntegrations() ([]*model.Integration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	integrations := make([]*model.Integration, 0, len(m.integrations))
	for _, integration := range m.integrations {
		integrations = append(integrations, integration)
	}

	return integrations, nil
}
