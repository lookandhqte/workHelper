package database

import (
	"errors"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"gorm.io/gorm"
)

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
