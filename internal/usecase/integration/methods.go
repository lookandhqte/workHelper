package integration

import (
	"context"
	"sync"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/dto"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

func (uc *IntegrationUseCase) Create(integration *entity.Integration) error {
	return uc.repo.AddIntegration(integration)
}

func (uc *IntegrationUseCase) Return(integration *entity.Integration) ([]*entity.Integration, error) {
	return uc.repo.GetIntegrations()
}

func (uc *IntegrationUseCase) Update(integration *entity.Integration) error {
	return uc.repo.UpdateIntegration(integration)
}

func (uc *IntegrationUseCase) Start(wg *sync.WaitGroup) func(ctx context.Context) {
	return uc.repo.Start(wg)
}

func (uc *IntegrationUseCase) GetIntegration(id int) (*entity.Integration, error) {
	return uc.repo.GetIntegration(id)
}

func (uc *IntegrationUseCase) GetIntegrationByClientID(client_id string) (*entity.Integration, error) {
	return uc.repo.GetIntegrationByClientID(client_id)
}

func (uc *IntegrationUseCase) GetContacts(token string) (*dto.ContactsResponse, error) {
	return uc.repo.GetContacts(token)
}

func (uc *IntegrationUseCase) Delete(id int) error {
	return uc.repo.DeleteIntegration(id)
}
