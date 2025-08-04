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
	GetAccountIntegrations(accountID int) (*entity.Integration, error)
	GetActiveAccount() *entity.Account
	ChangeActiveAccount(new_id int) error
	GetAccountWithCache(id int) (*entity.Account, error)
	UpdateAccount(account *entity.Account) error
	DeleteAccount(id int) error
	Exists(obj interface{}) bool
}

func New(r accountRepo) *AccountUseCase {
	return &AccountUseCase{repo: r}
}
