package entity

import (
	"time"
)

//Структура аккаунта
type Account struct {
	ID                    int       `json:"id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresIn  int       `json:"access_token_expires_in"`
	RefreshTokenExpiresIn int       `json:"refresh_token_expires_in"`
	CacheExpires          int       `json:"cache_expires"`
	CreatedAt             time.Time `json:"created_at"`
}

type Token struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ServerTime   int    `json:"server_time"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

//Структура интеграции
type Integration struct {
	AccountID   int    `json:"account_id"`
	SecretKey   string `json:"secret_key"`
	ClientID    string `json:"client_id"`
	RedirectUrl string `json:"redirect_url"`
	AuthCode    string `json:"auth_code"`
}
