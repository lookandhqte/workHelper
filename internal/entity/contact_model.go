package entity

//Contact структура контакта
type Contact struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AccountID int    `json:"account_id" gorm:"foreignKey:AccountID"`
	Status    string `json:"status"`
}
