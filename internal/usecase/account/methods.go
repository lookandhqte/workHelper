package account

import "github.com/lookandhqte/workHelper/internal/entity"

// Create создает аккаунт
func (uc *UseCase) Create(account *entity.Account) error {
	return uc.repo.AddAccount(account)
}

// ReturnOne возвращает аккаунт из хранилища
func (uc *UseCase) Return() (*entity.Account, error) {
	return uc.repo.GetAccount()
}

// Update обновляет аккаунт в хранилище
func (uc *UseCase) Update(account *entity.Account) error {
	return uc.repo.UpdateAccount(account)
}

// Delete удаляет аккаунт в хранилище
func (uc *UseCase) Delete() error {
	return uc.repo.DeleteAccount()
}
