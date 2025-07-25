package repository

import (
	"amocrm_golang/model"
)

//Определение методов для аккаунтов
type AccountRepository interface {
	AddAccount(account *model.Account) error
	GetAccounts() ([]*model.Account, error)
	UpdateAccount(account *model.Account) error
	DeleteAccount(id int) error
	GetAccount(id int) (*model.Account, error)
}

//Определение методов для интеграций
type IntegrationRepository interface {
	AddIntegration(integration *model.Integration) error
	GetIntegrations() ([]*model.Integration, error)
	UpdateIntegration(integration *model.Integration) error
	DeleteIntegration(accountID int) error
	GetAccountIntegrations(accountID int) (*model.Integration, error)
}

//Определение общего репоитория для аккаунтов и интеграций
type Repository interface {
	AccountRepository
	IntegrationRepository
}
