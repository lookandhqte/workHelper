package entity

//Token структура токена
type Token struct {
	TokenID      int    `json:"token_id" gorm:"primaryKey"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ServerTime   int    `json:"server_time"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UnisenderKey string `json:"unisender_key"`
}
