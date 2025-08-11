package storage

import "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"

//Storage определяет методы хранилища
type Storage interface {
	AddAccount(account *entity.Account) error
	GetAccounts() (*[]entity.Account, error)
	GetAccount(id int) (*entity.Account, error)
	GetAccountIntegrations(accountID int) (*[]entity.Integration, error)
	UpdateAccount(account *entity.Account) error
	DeleteAccount(id int) error
	AddIntegration(integration *entity.Integration) error
	GetIntegration(id int) (*entity.Integration, error)
	GetIntegrations() (*[]entity.Integration, error)
	UpdateIntegration(integration *entity.Integration) error
	DeleteIntegration(accountID int) error
	ReturnByClientID(clientID string) (*entity.Integration, error)
}

//DB Глобальная переенная хранилища
var DB Storage
