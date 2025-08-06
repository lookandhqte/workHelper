package persistent

import (
	"errors"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseStorage struct {
	DB *gorm.DB
}

func NewDatabaseStorage(cfg *config.Config) (*DatabaseStorage, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(
		&entity.Account{},
		&entity.Integration{},
		&entity.Contact{},
	)
	if err != nil {
		return nil, err
	}

	return &DatabaseStorage{DB: db}, nil
}

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

// Методы для интеграций

func (d *DatabaseStorage) AddIntegration(integration *entity.Integration) error {
	result := d.DB.Create(integration)
	return result.Error
}

func (d *DatabaseStorage) GetIntegration(id int) (*entity.Integration, error) {
	var integration entity.Integration
	result := d.DB.First(&integration, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("integration not found")
	}
	return &integration, result.Error
}

func (d *DatabaseStorage) GetIntegrations() (*[]entity.Integration, error) {
	var integrations []entity.Integration
	result := d.DB.Find(&integrations)
	return &integrations, result.Error
}

func (d *DatabaseStorage) UpdateIntegration(integration *entity.Integration) error {
	result := d.DB.Save(integration)
	return result.Error
}

func (d *DatabaseStorage) DeleteIntegration(id int) error {
	result := d.DB.Delete(&entity.Integration{}, id)
	return result.Error
}

func (d *DatabaseStorage) ReturnByClientID(clientID string) (*entity.Integration, error) {
	var integration *entity.Integration
	result := d.DB.Where("client_id = ?", clientID).First(&integration)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("integration not found")
	}
	return integration, result.Error
}
