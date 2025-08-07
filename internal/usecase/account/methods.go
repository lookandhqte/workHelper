package account

import "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"

//Create создает аккаунт
func (uc *UseCase) Create(account *entity.Account) error {
	return uc.repo.AddAccount(account)
}

//ReturnAll возвращает все аккаунты из хранилища
func (uc *UseCase) ReturnAll() ([]*entity.Account, error) {
	return uc.repo.GetAccounts()
}

//ReturnOne возвращает аккаунт из хранилища
func (uc *UseCase) ReturnOne(id int) (*entity.Account, error) {
	return uc.repo.GetAccount(id)
}

//SaveContacts сохраняет контакты
func (uc *UseCase) SaveContacts(contact *[]entity.Contact) error {
	return uc.repo.SaveContacts(contact)
}

//ReturnIntegrations возвращает интеграции аккаунта из хранилища
func (uc *UseCase) ReturnIntegrations(accountID int) (*[]entity.Integration, error) {
	return uc.repo.GetAccountIntegrations(accountID)
}

//Update обновляет аккаунт в хранилище
func (uc *UseCase) Update(account *entity.Account) error {
	return uc.repo.UpdateAccount(account)
}

//Delete удаляет аккаунт в хранилище
func (uc *UseCase) Delete(id int) error {
	return uc.repo.DeleteAccount(id)
}
