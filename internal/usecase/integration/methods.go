package integration

import "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"

func (uc *IntegrationUseCase) Create(integration *entity.Integration) error {
	return uc.repo.AddIntegration(integration)
}

func (uc *IntegrationUseCase) Return(integration *entity.Integration) ([]*entity.Integration, error) {
	return uc.repo.GetIntegrations()
}

func (uc *IntegrationUseCase) Update(integration *entity.Integration) error {
	return uc.repo.UpdateIntegration(integration)
}

func (uc *IntegrationUseCase) UpdateTokens(tokens *entity.Token) error {
	return uc.repo.UpdateTokens(tokens)
}

func (uc *IntegrationUseCase) GetIntegration(id int) (*entity.Integration, error) {
	return uc.repo.GetIntegration(id)
}

func (uc *IntegrationUseCase) Exists(obj interface{}) bool {
	return uc.repo.Exists(obj)
}

func (uc *IntegrationUseCase) GetTokensByAuthCode(code string, client_id string) (*entity.Token, error) {
	return uc.repo.GetTokensByAuthCode(code, client_id)
}

func (uc *IntegrationUseCase) GetActiveIntegrations() ([]*entity.Integration, error) {
	return uc.repo.GetActiveIntegrations()
}

func (uc *IntegrationUseCase) Delete(id int) error {
	return uc.repo.DeleteIntegration(id)
}

func (uc *IntegrationUseCase) GetIntegrationByClientID(client_id string) (*entity.Integration, error) {
	return uc.repo.GetIntegrationByClientID(client_id)
}
func (uc *IntegrationUseCase) MakeIntegrationActive(new_id int) error {
	return uc.repo.MakeIntegrationActive(new_id)
}

func (uc *IntegrationUseCase) MakeIntegrationInactive(new_id int) error {
	return uc.repo.MakeIntegrationInactive(new_id)
}
