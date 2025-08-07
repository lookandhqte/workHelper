package entity

//Integration структура интеграции
type Integration struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	AccountID   int    `json:"account_id" gorm:"foreignKey:Account.ID"`
	SecretKey   string `json:"secret_key"`
	ClientID    string `json:"client_id"`
	RedirectURL string `json:"redirect_url"`
	AuthCode    string `json:"auth_code"`
	TokenID     int    `json:"integration_tokens" gorm:"foreignKey:Token.ID"`
}
