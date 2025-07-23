package model

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID `json:"id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expires      time.Time `json:"expires"`
}

type Integration struct {
	AccountID   uuid.UUID `json:"account_id"`
	SecretKey   string    `json:"secret_key"`
	ClientID    string    `json:"client_id"`
	RedirectUrl string    `json:"redirect_url"`
	AuthCode    string    `json:"auth_code"`
}

type Repository struct {
	db *Database
}

type Database struct {
}
