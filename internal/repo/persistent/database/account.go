package database

import (
	"errors"

	"github.com/lookandhqte/workHelper/internal/entity"
	"gorm.io/gorm"
)

// AddAccount создает или заменяет аккаунт (всегда только один)
func (d *Storage) AddAccount(account *entity.Account) error {
	d.AddToken(&account.Token)
	return d.DB.Where("1 = 1").Delete(&entity.Account{}).Create(account).Error
}

// GetAccount возвращает единственный аккаунт
func (d *Storage) GetAccount() (*entity.Account, error) {
	var account entity.Account
	result := d.DB.Preload("Token").First(&account)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("account not found")
	}
	return &account, result.Error
}

// UpdateAccount обновляет единственный аккаунт
func (d *Storage) UpdateAccount(account *entity.Account) error {
	var existingAccount entity.Account
	result := d.DB.First(&existingAccount)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("account not found, cannot update")
	}
	if existingAccount.Token != account.Token {
		d.AddToken(&account.Token)
	}
	return d.DB.Model(&existingAccount).Updates(account).Error
}

// DeleteAccount удаляет единственный аккаунт
func (d *Storage) DeleteAccount() error {
	return d.DB.Where("1 = 1").Delete(&entity.Account{}).Error
}
