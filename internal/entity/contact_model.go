package entity

//Contact структура контакта
type Contact struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	AccountID int    `json:"account_id" gorm:"foreignKey:AccountID;references:ID"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
}
