package integration

import (
	"amocrm_golang/internal/entity"
)

type IntegrationUseCase struct {
	repo integrationRepo
}

type integrationRepo interface {
	AddIntegration(integration *entity.Integration) error
	GetIntegrations() ([]*entity.Integration, error)
	UpdateIntegration(integration *entity.Integration) error
	DeleteIntegration(accountID int) error
}

func New(r integrationRepo) *IntegrationUseCase {
	return &IntegrationUseCase{repo: r}
}
