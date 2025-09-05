package entity

// Token структура токена
type Token struct {
	AccountID    int    `json:"account_id" gorm:"primaryKey, foreignKey:Account.ID"`
	CreatedAt    int    `json:"created_at"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
