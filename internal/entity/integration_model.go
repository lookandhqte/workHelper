package entity

//Integration структура интеграции
type Integration struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	AccountID   int    `json:"account_id" gorm:"foreignKey:AccountID;references:ID"`
	SecretKey   string `json:"secret_key"`
	ClientID    string `json:"client_id"`
	RedirectURL string `json:"redirect_url"`
	AuthCode    string `json:"auth_code"`
	Token       *Token `json:"integration_token"`
}
