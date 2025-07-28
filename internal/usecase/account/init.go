package account

import (
	"amocrm_golang/internal/entity"
	"amocrm_golang/internal/repo"
)

type AccountUseCase struct {
	repo repo.AccountRepository
}

type AccountRepo interface {
	AddAccount(account *entity.Account) error
	GetAccounts() ([]*entity.Account, error)
	GetAccount(id int) (*entity.Account, error)
	GetAccountIntegrations(accountID int) (*entity.Integration, error)
	UpdateAccount(account *entity.Account) error
	DeleteAccount(id int) error
}

func New(r AccountRepo) *AccountUseCase {
	return &AccountUseCase{repo: r}
}
