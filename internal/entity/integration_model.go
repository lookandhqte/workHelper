package entity

//Структура интеграции
type Integration struct {
	AccountID   int    `json:"account_id" gorm:"primaryKey"`
	SecretKey   string `json:"secret_key"`
	ClientID    string `json:"client_id"`
	RedirectUrl string `json:"redirect_url"`
	AuthCode    string `json:"auth_code"`
	Token       *Token `json:"integration_tokens"`
}
