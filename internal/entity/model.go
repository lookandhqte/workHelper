package entity

//Структура аккаунта
type Account struct {
	ID                    int           `json:"id"`
	AccessToken           string        `json:"access_token"`
	RefreshToken          string        `json:"refresh_token"`
	AccessTokenExpiresIn  int           `json:"access_token_expires_in"`
	RefreshTokenExpiresIn int           `json:"refresh_token_expires_in"`
	CacheExpires          int           `json:"cache_expires"`
	CreatedAt             int           `json:"created_at"`
	Integrations          *Integrations `json:"integrations"`
}

type Token struct {
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
	ServerTime    int    `json:"server_time"`
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	IntegrationID int    `json:"integration_id"`
}

type Contact struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Contacts []Contact

type Integrations []Integration

//Структура интеграции
type Integration struct {
	AccountID       int       `json:"account_id"`
	SecretKey       string    `json:"secret_key"`
	ClientID        string    `json:"client_id"`
	RedirectUrl     string    `json:"redirect_url"`
	AuthCode        string    `json:"auth_code"`
	Token           *Token    `json:"integration_tokens"`
	AccountContacts *Contacts `json:"account_contacts"`
}
