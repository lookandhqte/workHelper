package account

import (
	entity "github.com/lookandhqte/workHelper/internal/entity"
)

// UseCase структура
type UseCase struct {
	repo accountRepo
}

// accountRepo абстракция для определения методов репозитория
type accountRepo interface {
	AddAccount(account *entity.Account) error
	GetAccount() (*entity.Account, error)
	UpdateAccount(account *entity.Account) error
	DeleteAccount() error
}

// New создает новый репозиторий
func New(r accountRepo) *UseCase {
	return &UseCase{repo: r}
}
