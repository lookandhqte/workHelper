package repository

import (
	"amocrm_golang/model"
)

type AccountRepository interface {
	AddAccount(account *model.Account) error
	GetAccounts() ([]*model.Account, error)
	UpdateAccount(account *model.Account) error
	DeleteAccount(id int) error
	GetAccount(id int) (*model.Account, error)
}

type IntegrationRepository interface {
	AddIntegration(integration *model.Integration) error
	GetIntegrations() ([]*model.Integration, error)
	UpdateIntegration(integration *model.Integration) error
	DeleteIntegration(accountID int) error
	GetAccountIntegrations(accountID int) (*model.Integration, error)
}

type Repository interface {
	AccountRepository
	IntegrationRepository
}
