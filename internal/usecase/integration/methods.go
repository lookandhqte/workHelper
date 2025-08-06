package integration

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

func (uc *IntegrationUseCase) Create(integration *entity.Integration) error {
	return uc.repo.AddIntegration(integration)
}

func (uc *IntegrationUseCase) Return() (*[]entity.Integration, error) {
	return uc.repo.GetIntegrations()
}

func (uc *IntegrationUseCase) ReturnByClientID(client_id string) (*entity.Integration, error) {
	return uc.repo.ReturnByClientID(client_id)
}

func (uc *IntegrationUseCase) Update(integration *entity.Integration) error {
	return uc.repo.UpdateIntegration(integration)
}

func (uc *IntegrationUseCase) GetIntegration(id int) (*entity.Integration, error) {
	return uc.repo.GetIntegration(id)
}

func (uc *IntegrationUseCase) Delete(id int) error {
	return uc.repo.DeleteIntegration(id)
}
