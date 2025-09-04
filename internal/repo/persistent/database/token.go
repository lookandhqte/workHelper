package database

import (
	"github.com/lookandhqte/workHelper/internal/entity"
)

// AddToken создает или заменяет токен (всегда только один)
func (d *Storage) AddToken(token *entity.Token) error {
	return d.DB.Where("1 = 1").Delete(&entity.Token{}).Create(token).Error
}

// GetTokenExpiry возвращает время в секундах, когда токен истечет
func (d *Storage) GetTokenExpiry() (int, error) {
	var token entity.Token
	result := d.DB.First(&token)
	return token.ExpiresIn, result.Error
}
