package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

//Структура аккаунта
type Account struct {
	ID              int           `json:"id" gorm:"primaryKey"`
	CacheExpires    int           `json:"cache_expires"`
	CreatedAt       int           `json:"created_at"`
	Integrations    []Integration `json:"integrations" gorm:"foreignKey:AccountID"`
	AccountContacts []Contact     `json:"contacts" gorm:"foreignKey:AccountID"`
}

type Token struct {
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
	ServerTime    int    `json:"server_time"`
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	IntegrationID int    `json:"integration_id"`
}

func (t *Token) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *Token) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type Contact struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AccountID int    `json:"account_id"`
}

//Структура интеграции
type Integration struct {
	AccountID   int    `json:"account_id" gorm:"primaryKey"`
	SecretKey   string `json:"secret_key"`
	ClientID    string `json:"client_id"`
	RedirectUrl string `json:"redirect_url"`
	AuthCode    string `json:"auth_code"`
	Token       *Token `json:"integration_tokens"`
}
