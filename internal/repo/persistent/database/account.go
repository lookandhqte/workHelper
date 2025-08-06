package database

import (
	"errors"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"gorm.io/gorm"
)

func (d *DatabaseStorage) AddAccount(account *entity.Account) error {
	result := d.DB.Create(account)
	return result.Error
}

func (d *DatabaseStorage) GetAccounts() ([]*entity.Account, error) {
	var accounts []*entity.Account
	result := d.DB.Find(&accounts)
	return accounts, result.Error
}

func (d *DatabaseStorage) GetAccount(id int) (*entity.Account, error) {
	var account entity.Account
	result := d.DB.First(&account, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("account not found")
	}
	return &account, result.Error
}

func (d *DatabaseStorage) GetAccountIntegrations(accountID int) (*[]entity.Integration, error) {
	var integrations []entity.Integration
	result := d.DB.Where("account_id = ?", accountID).Find(&integrations)
	return &integrations, result.Error
}

func (d *DatabaseStorage) UpdateAccount(account *entity.Account) error {
	result := d.DB.Save(account)
	return result.Error
}

func (d *DatabaseStorage) DeleteAccount(id int) error {
	result := d.DB.Delete(&entity.Account{}, id)
	return result.Error
}
