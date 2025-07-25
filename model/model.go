package model

import (
	"time"
)

type Account struct {
	ID           int       `json:"id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expires      int       `json:"expires"`
	CreatedAt    time.Time `json:"created_at"`
	TokenExpires time.Time `json:"token_expires"`
}

type Integration struct {
	AccountID   int    `json:"account_id"`
	SecretKey   string `json:"secret_key"`
	ClientID    string `json:"client_id"`
	RedirectUrl string `json:"redirect_url"`
	AuthCode    string `json:"auth_code"`
}

type Repository struct {
	db *Database
}

type Database struct {
}
