package repo

import (
	"amocrm_golang/internal/entity"
)

//Определение методов для аккаунтов
type AccountRepository interface {
	AddAccount(account *entity.Account) error
	GetAccounts() ([]*entity.Account, error)
	UpdateAccount(account *entity.Account) error
	DeleteAccount(id int) error
	GetAccount(id int) (*entity.Account, error)
	GetAccountIntegrations(accountID int) (*entity.Integration, error)
}

//Определение методов для интеграций
type IntegrationRepository interface {
	AddIntegration(integration *entity.Integration) error
	GetIntegrations() ([]*entity.Integration, error)
	UpdateIntegration(integration *entity.Integration) error
	DeleteIntegration(accountID int) error
}

//Определение общего репоитория для аккаунтов и интеграций
type Repository interface {
	AccountRepository
	IntegrationRepository
}
