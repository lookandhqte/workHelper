package database

import (
	"errors"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"gorm.io/gorm"
)

//AddIntegration добавляет интеграцию
func (d *Storage) AddIntegration(integration *entity.Integration) error {
	result := d.DB.Create(integration)
	return result.Error
}

//GetIntegration возвращает интеграцию по id
func (d *Storage) GetIntegration(id int) (*entity.Integration, error) {
	var integration entity.Integration
	result := d.DB.First(&integration, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("integration not found")
	}
	return &integration, result.Error
}

//GetIntegrations возвращает все интеграции
func (d *Storage) GetIntegrations() (*[]entity.Integration, error) {
	var integrations []entity.Integration
	result := d.DB.Find(&integrations)
	return &integrations, result.Error
}

//UpdateIntegration обновляет интеграцию
func (d *Storage) UpdateIntegration(integration *entity.Integration) error {
	result := d.DB.Save(integration)
	return result.Error
}

//DeleteIntegration возвращает интеграцию по clientID
func (d *Storage) DeleteIntegration(id int) error {
	result := d.DB.Delete(&entity.Integration{}, id)
	return result.Error
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
