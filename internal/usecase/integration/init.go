package integration

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

type IntegrationUseCase struct {
	repo integrationRepo
}

type integrationRepo interface {
	AddIntegration(integration *entity.Integration) error
	GetIntegrations() (*[]entity.Integration, error)
	GetIntegration(id int) (*entity.Integration, error)
	UpdateIntegration(integration *entity.Integration) error
	DeleteIntegration(accountID int) error
	ReturnByClientID(client_id string) (*entity.Integration, error)
}

func New(r integrationRepo) *IntegrationUseCase {
	return &IntegrationUseCase{repo: r}
}
