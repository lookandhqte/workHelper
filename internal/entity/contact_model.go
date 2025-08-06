package entity

//Contact структура контакта
type Contact struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AccountID int    `json:"account_id" gorm:"primaryKey"`
}
