package entity

import (
	"time"
)

//Структура аккаунта
type Account struct {
	ID           int          `json:"id"`
	AccessToken  AccessToken  `json:"access_token"`
	RefreshToken RefreshToken `json:"refresh_token"`
	CacheExpires int          `json:"cache_expires"`
	CreatedAt    time.Time    `json:"created_at"`
}

type Token struct {
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
	ServerTime   int          `json:"server_time"`
	AccessToken  AccessToken  `json:"access_token"`
	RefreshToken RefreshToken `json:"refresh_token"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	CreatedAt   int    `json:"created_at"`
	TTL         int    `json:"ttl"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	CreatedAt    int    `json:"created_at"`
	TTL          int    `json:"ttl"`
}

//Структура интеграции
type Integration struct {
	AccountID   int    `json:"account_id"`
	SecretKey   string `json:"secret_key"`
	ClientID    string `json:"client_id"`
	RedirectUrl string `json:"redirect_url"`
	AuthCode    string `json:"auth_code"`
}
