package storage

import "github.com/lookandhqte/workHelper/internal/entity"

// Storage определяет методы хранилища
type Storage interface {
	AddAccount(account *entity.Account) error
	GetAccount() (*entity.Account, error)
	UpdateAccount(account *entity.Account) error
	DeleteAccount() error
}

// DB Глобальная переенная хранилища
var DB Storage
