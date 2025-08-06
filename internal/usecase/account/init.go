package account

import (
	entity "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

//UseCase структура
type UseCase struct {
	repo accountRepo
}

//accountRepo абстракция для определения методов репозитория
type accountRepo interface {
	AddAccount(account *entity.Account) error
	GetAccounts() ([]*entity.Account, error)
	GetAccount(id int) (*entity.Account, error)
	GetAccountIntegrations(accountID int) (*[]entity.Integration, error)
	UpdateAccount(account *entity.Account) error
	DeleteAccount(id int) error
}

//New создает новый репозиторий
func New(r accountRepo) *UseCase {
	return &UseCase{repo: r}
}
