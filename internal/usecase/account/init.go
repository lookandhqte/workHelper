package account

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

type AccountUseCase struct {
	repo accountRepo
}

type accountRepo interface {
	AddAccount(account *entity.Account) error
	GetAccounts() ([]*entity.Account, error)
	GetAccount(id int) (*entity.Account, error)
	ChangeActiveAccount(new_id int) error
	GetAccountIntegrations(accountID int) (*entity.Integration, error)
	UpdateAccount(account *entity.Account) error
	DeleteAccount(id int) error
}

func New(r accountRepo) *AccountUseCase {
	return &AccountUseCase{repo: r}
}
