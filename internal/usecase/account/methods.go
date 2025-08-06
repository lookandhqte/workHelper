package account

import "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"

func (uc *AccountUseCase) Create(account *entity.Account) error {
	return uc.repo.AddAccount(account)
}

func (uc *AccountUseCase) GetAccounts() ([]*entity.Account, error) {
	return uc.repo.GetAccounts()
}

func (uc *AccountUseCase) GetAccount(id int) (*entity.Account, error) {
	return uc.repo.GetAccount(id)
}

func (uc *AccountUseCase) GetAccountIntegrations(accountID int) (*[]entity.Integration, error) {
	return uc.repo.GetAccountIntegrations(accountID)
}

func (uc *AccountUseCase) Update(account *entity.Account) error {
	return uc.repo.UpdateAccount(account)
}

func (uc *AccountUseCase) Delete(id int) error {
	return uc.repo.DeleteAccount(id)
}
