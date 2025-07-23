package repository

import (
	"amocrm_golang/model"

	"github.com/google/uuid"
)

type AccountRepository interface {
	AddAccount(account *model.Account) error
	GetAccounts() ([]*model.Account, error)
	UpdateAccount(account *model.Account) error
	DeleteAccount(id uuid.UUID) error
	GetAccount(id uuid.UUID) (*model.Account, error)
}

type IntegrationRepository interface {
	AddIntegration(integration *model.Integration) error
	GetIntegrations() ([]*model.Integration, error)
	UpdateIntegration(integration *model.Integration) error
	DeleteIntegration(accountID uuid.UUID) error
	GetAccountIntegrations(accountID uuid.UUID) (*model.Integration, error)
}

type Repository interface {
	AccountRepository
	IntegrationRepository
}
