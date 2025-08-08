package entity

//Token структура токена
type Token struct {
	IntegrationID int    `json:"int_id" gorm:"primaryKey, foreignKey:IntegrationID;references:ID"`
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
	ServerTime    int    `json:"server_time"`
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	UnisenderKey  string `json:"unisender_key"`
}
