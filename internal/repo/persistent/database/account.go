package database

import (
	"errors"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"gorm.io/gorm"
)

//AddAccount создает аккаунт
func (d *Storage) AddAccount(account *entity.Account) error {
	result := d.DB.Create(account)
	return result.Error
}

//GetAccounts возвращает все аккаунты
func (d *Storage) GetAccounts() ([]*entity.Account, error) {
	var accounts []*entity.Account
	result := d.DB.Find(&accounts)
	return accounts, result.Error
}

//GetAccount возвращает аккаунт
func (d *Storage) GetAccount(id int) (*entity.Account, error) {
	var account entity.Account
	result := d.DB.First(&account, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("account not found")
	}
	return &account, result.Error
}

//GetAccountIntegrations возвращает все интеграции аккаунта
func (d *Storage) GetAccountIntegrations(accountID int) (*[]entity.Integration, error) {
	var integrations []entity.Integration
	result := d.DB.Where("account_id = ?", accountID).Find(&integrations)
	return &integrations, result.Error
}

//UpdateAccount обновляет аккаунт
func (d *Storage) UpdateAccount(account *entity.Account) error {
	result := d.DB.Save(account)
	return result.Error
}

//DeleteAccount удаляет аккаунт
func (d *Storage) DeleteAccount(id int) error {
	result := d.DB.Delete(&entity.Account{}, id)
	return result.Error
}
