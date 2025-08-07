package entity

//Token структура токена
type Token struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	AccountID    int    `json:"tok_account_id" gorm:"foreignKey: Account.ID"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ServerTime   int    `json:"server_time"`   //
	AccessToken  string `json:"access_token"`  //
	RefreshToken string `json:"refresh_token"` //
	UnisenderKey string `json:"unisender_key"` //
}
