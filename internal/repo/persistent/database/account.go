package database

import (
	"errors"
	"time"

	"github.com/lookandhqte/workHelper/internal/entity"
	"gorm.io/gorm"
)

// AddAccount создает или заменяет аккаунт (всегда только один)
func (d *Storage) AddAccount(account *entity.Account) error {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("1 = 1").Delete(&entity.Token{}).Error; err != nil {
			return err
		}
		if err := tx.Where("1 = 1").Delete(&entity.Account{}).Error; err != nil {
			return err
		}
		account.CreatedAt = int(time.Now().Unix())

		if err := tx.Create(account).Error; err != nil {
			return err
		}
		if account.Token != (entity.Token{}) {
			account.Token.AccountID = account.ID
			account.Token.CreatedAt = int(time.Now().Unix())
			if err := tx.Create(&account.Token).Error; err != nil {
				return err
			}
		}

		return nil
	})
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
	return d.DB.Transaction(func(tx *gorm.DB) error {
		var existingAccount entity.Account
		if err := tx.First(&existingAccount).Error; err != nil {
			return errors.New("account not found, cannot update")
		}

		accountCopy := *account
		accountCopy.Token = entity.Token{}
		if err := tx.Model(&existingAccount).Updates(accountCopy).Error; err != nil {
			return err
		}

		if account.Token != (entity.Token{}) {
			account.Token.AccountID = existingAccount.ID
			account.Token.CreatedAt = int(time.Now().Unix())
			if err := tx.Where("1 = 1").Delete(&entity.Token{}).Error; err != nil {
				return err
			}
			if err := tx.Create(&account.Token).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteAccount удаляет единственный аккаунт
func (d *Storage) DeleteAccount() error {
	return d.DB.Where("1 = 1").Delete(&entity.Account{}).Error
}
