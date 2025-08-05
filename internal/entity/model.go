package entity

//Структура аккаунта
type Account struct {
	ID           int            `json:"id"`
	CacheExpires int            `json:"cache_expires"`
	CreatedAt    int            `json:"created_at"`
	Integrations *[]Integration `json:"integrations"`
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

//Структура интеграции
type Integration struct {
	AccountID       int        `json:"account_id"`
	SecretKey       string     `json:"secret_key"`
	ClientID        string     `json:"client_id"`
	RedirectUrl     string     `json:"redirect_url"`
	AuthCode        string     `json:"auth_code"`
	Token           *Token     `json:"integration_tokens"`
	AccountContacts *[]Contact `json:"account_contacts"`
}
