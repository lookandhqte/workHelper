package integration

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

type IntegrationUseCase struct {
	repo integrationRepo
}

type integrationRepo interface {
	AddIntegration(integration *entity.Integration) error
	GetIntegrations() ([]*entity.Integration, error)
	GetIntegration(id int) (*entity.Integration, error)
	UpdateIntegration(integration *entity.Integration) error
	GetTokensByAuthCode(code string, client_id string) (*entity.Token, error)
	UpdateTokens(tokens *entity.Token) error
	AddTokens(tokens *entity.Token) error
	GetActiveIntegrations() ([]*entity.Integration, error)
	DeleteIntegration(accountID int) error
	DeleteTokens() error
	GetIntegrationByClientID(client_id string) (*entity.Integration, error)
	MakeIntegrationActive(new_id int) error
	MakeIntegrationInactive(new_id int) error
	GetRefreshToken() (string, error)
	GetTokens() (*entity.Token, error)
}

func New(r integrationRepo) *IntegrationUseCase {
	return &IntegrationUseCase{repo: r}
}
