package database

import (
	"errors"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"gorm.io/gorm"
)

//AddIntegration добавляет интеграцию
func (d *Storage) AddIntegration(integration *entity.Integration) error {
	var tokens *entity.Token = &entity.Token{}
	tokens.AccountID = integration.AccountID
	integration.Token = *tokens
	result := d.DB.Create(integration)
	return result.Error
}

//GetIntegration возвращает интеграцию по id
func (d *Storage) GetIntegration(id int) (*entity.Integration, error) {
	var integration *entity.Integration
	result := d.DB.First(&integration, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("integration not found")
	}
	return integration, result.Error
}

//GetIntegrations возвращает все интеграции
func (d *Storage) GetIntegrations() (*[]entity.Integration, error) {
	var integrations []entity.Integration
	integrationsPtr := &integrations
	result := d.DB.Find(&integrations)
	return integrationsPtr, result.Error
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

//UpdateToken обновляет токены
func (d *Storage) UpdateToken(token *entity.Token) error {
	if err := d.DB.Where("account_id = ?", token.AccountID).Save(token).Error; err != nil {
		return err
	}
	return nil
}

//GetTokens возвращает токены интеграции
func (d *Storage) GetTokens(id int) (*entity.Token, error) {
	var token *entity.Token
	result := d.DB.First(&token, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("integration token not found")
	}
	return token, result.Error
}
