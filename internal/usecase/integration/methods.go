package integration

import (
	entity "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

//Create создает новую интеграцию в хранилище
func (uc *UseCase) Create(integration *entity.Integration) error {
	return uc.repo.AddIntegration(integration)
}

//ReturnAll возвращает все интеграции из хранилища
func (uc *UseCase) ReturnAll() (*[]entity.Integration, error) {
	return uc.repo.GetIntegrations()
}

//ReturnByClientID возвращает интеграцию по параметру client_id
func (uc *UseCase) ReturnByClientID(clientID string) (*entity.Integration, error) {
	return uc.repo.ReturnByClientID(clientID)
}

//Update обновляет интеграцию
func (uc *UseCase) Update(integration *entity.Integration) error {
	return uc.repo.UpdateIntegration(integration)
}

//ReturnOne возвращает интеграцию по id
func (uc *UseCase) ReturnOne(id int) (*entity.Integration, error) {
	return uc.repo.GetIntegration(id)
}

//Delete удаляет интеграцию
func (uc *UseCase) Delete(id int) error {
	return uc.repo.DeleteIntegration(id)
}

//UpdateToken обновляет токены
func (uc *UseCase) UpdateToken(token *entity.Token) error {
	return uc.repo.UpdateToken(token)
}

func (uc *UseCase) GetTokens(id int) (*entity.Token, error) {
	return uc.repo.GetTokens(id)
}
