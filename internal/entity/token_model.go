package entity

//Token структура токена
type Token struct {
	AccountID     int    `json:"tok_account_id" gorm:"primaryKey, foreignKey: Account.ID"`
	IntegrationID int    `json:"int_id" gorm:"foreignKey: Integration.ID"`
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
	ServerTime    int    `json:"server_time"`
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	UnisenderKey  string `json:"unisender_key"`
}
