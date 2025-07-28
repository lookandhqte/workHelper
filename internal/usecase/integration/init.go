package integration

import (
	"amocrm_golang/internal/entity"
	"amocrm_golang/internal/repo"
)

type IntegrationUseCase struct {
	repo repo.IntegrationRepository
}

type IntegrationRepo interface {
	AddIntegration(integration *entity.Integration) error
	GetIntegrations() ([]*entity.Integration, error)
	UpdateIntegration(integration *entity.Integration) error
	DeleteIntegration(accountID int) error
}

func New(r IntegrationRepo) *IntegrationUseCase {
	return &IntegrationUseCase{repo: r}
}
