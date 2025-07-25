package usecase

import (
	"amocrm_golang/internal/entity"
	"amocrm_golang/internal/repo"
)

type IntegrationUseCase struct {
	repo repo.IntegrationRepository
}

func NewIntegrationUseCase(r repo.IntegrationRepository) *IntegrationUseCase {
	return &IntegrationUseCase{r}
}

func (uc *IntegrationUseCase) Create(integration *entity.Integration) error {
	return uc.repo.AddIntegration(integration)
}

func (uc *IntegrationUseCase) Return(integration *entity.Integration) ([]*entity.Integration, error) {
	return uc.repo.GetIntegrations()
}

func (uc *IntegrationUseCase) Update(integration *entity.Integration) error {
	return uc.repo.UpdateIntegration(integration)
}

func (uc *IntegrationUseCase) Delete(id int) error {
	return uc.repo.DeleteIntegration(id)
}
