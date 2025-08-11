package database

import (
	"errors"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"gorm.io/gorm"
)

//AddIntegration добавляет интеграцию
func (d *Storage) AddIntegration(integration *entity.Integration) error {
	return d.DB.Create(integration).Error
}

//GetIntegration возвращает интеграцию по id
func (d *Storage) GetIntegration(id int) (*entity.Integration, error) {
	var integration *entity.Integration
	result := d.DB.Preload("Token").First(&integration, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("integration not found")
	}
	return integration, result.Error
}

//GetIntegrations возвращает все интеграции
func (d *Storage) GetIntegrations() (*[]entity.Integration, error) {
	var integrations []entity.Integration
	integrationsPtr := &integrations
	result := d.DB.Preload("Token").Find(&integrationsPtr)
	return integrationsPtr, result.Error
}

//UpdateIntegration обновляет интеграцию
func (d *Storage) UpdateIntegration(integration *entity.Integration) error {
	if integration.Token.AccessToken != "" {
		d.DB.Model(&integration).Association("Token").Replace(&integration.Token)
	}
	result := d.DB.Where("id = ?", integration.ID).Updates(integration)
	return result.Error
}

//DeleteIntegration возвращает интеграцию по clientID
func (d *Storage) DeleteIntegration(id int) error {
	return d.DB.Delete(&entity.Integration{}, id).Error
}

//ReturnByClientID возвращает интеграцию по clientID
func (d *Storage) ReturnByClientID(clientID string) (*entity.Integration, error) {
	var integration *entity.Integration
	result := d.DB.Where("client_id = ?", clientID).First(&integration)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("integration not found")
	}
	return integration, result.Error
}
