package database

import (
	"amocrm_golang/model"
	"amocrm_golang/repository"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

//Структура сервиса аккаунтов
type AccountService struct {
	repo repository.AccountRepository
}

//Структура сервиса интеграций
type IntegrationService struct {
	repo repository.IntegrationRepository
}

//Создание нового сервиса аккаунтов
func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

//Создание нового сервиса интеграций
func NewIntegrationService(repo repository.IntegrationRepository) *IntegrationService {
	return &IntegrationService{repo: repo}
}

//бляблябля блюблюблю
func (s *AccountService) CreateAccount(account *model.Account) error {
	return s.repo.AddAccount(account)
}

func (s *AccountService) GetAccountList() ([]*model.Account, error) {
	return s.repo.GetAccounts()
}

func (s *AccountService) GetAccountByID(id uuid.UUID) (*model.Account, error) {
	return s.repo.GetAccount(id)
}

func (s *AccountService) UpdateAccount(account *model.Account) error {
	return s.repo.UpdateAccount(account)
}

func (s *AccountService) DeleteAccount(id uuid.UUID) error {
	return s.repo.DeleteAccount(id)
}

func (s *IntegrationService) CreateIntegration(integration *model.Integration) error {
	return s.repo.AddIntegration(integration)
}

func (s *IntegrationService) GetIntegrationList() ([]*model.Integration, error) {
	return s.repo.GetIntegrations()
}

func (s *IntegrationService) GetAccountIntegrations(accountID uuid.UUID) (*model.Integration, error) {
	return s.repo.GetAccountIntegrations(accountID)
}

func (s *IntegrationService) UpdateIntegration(integration *model.Integration) error {
	return s.repo.UpdateIntegration(integration)
}

func (s *IntegrationService) DeleteIntegration(id uuid.UUID) error {
	return s.repo.DeleteIntegration(id)
}

type MemoryStorage struct {
	mu           sync.RWMutex
	accounts     map[uuid.UUID]*model.Account
	integrations map[uuid.UUID]*model.Integration
}

// NewMemoryStorage создает новое in-memory хранилище.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		accounts:     make(map[uuid.UUID]*model.Account),
		integrations: make(map[uuid.UUID]*model.Integration),
	}
}

func (m *MemoryStorage) UpdateAccount(account *model.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if account.ID == uuid.Nil {
		return fmt.Errorf("account ID is required")
	}

	if _, exists := m.accounts[account.ID]; !exists {
		return fmt.Errorf("аккаунт не найден")
	}

	m.accounts[account.ID] = account
	return nil
}
func (m *MemoryStorage) DeleteAccount(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if id == uuid.Nil {
		return fmt.Errorf("ID is required")
	}

	if _, exists := m.accounts[id]; !exists {
		return fmt.Errorf("account not found")
	}

	delete(m.accounts, id)
	return nil
}

func (m *MemoryStorage) DeleteIntegration(accountID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if accountID == uuid.Nil {
		return fmt.Errorf("account ID is required")
	}

	if _, exists := m.integrations[accountID]; !exists {
		return fmt.Errorf("integration not found")
	}

	delete(m.integrations, accountID)
	return nil
}

func (m *MemoryStorage) UpdateIntegration(integration *model.Integration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.integrations[integration.AccountID]; !exists {
		return fmt.Errorf("интеграция не найдена")
	}

	m.integrations[integration.AccountID] = integration
	return nil
}

func (m *MemoryStorage) AddAccount(account *model.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.accounts[account.ID] = account
	return nil
}
func (m *MemoryStorage) GetAccounts() ([]*model.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	accounts := make([]*model.Account, 0, len(m.accounts))
	for _, account := range m.accounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (m *MemoryStorage) GetAccount(id uuid.UUID) (*model.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var account model.Account
	if _, exists := m.accounts[id]; !exists {
		return nil, fmt.Errorf("аккаунт не найден по id")
	}
	account = *m.accounts[id]
	return &account, nil
}
func (m *MemoryStorage) GetAccountIntegrations(accountID uuid.UUID) (*model.Integration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if accountID == uuid.Nil {
		return nil, fmt.Errorf("account ID is required")
	}

	integration, exists := m.integrations[accountID]

	if !exists {
		return nil, fmt.Errorf("integration for account %s not found", accountID)
	}
	return integration, nil
}

func (m *MemoryStorage) AddIntegration(integration *model.Integration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.integrations[integration.AccountID] = integration
	return nil
}
func (m *MemoryStorage) GetIntegrations() ([]*model.Integration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	integrations := make([]*model.Integration, 0, len(m.integrations))
	for _, integration := range m.integrations {
		integrations = append(integrations, integration)
	}

	return integrations, nil
}
