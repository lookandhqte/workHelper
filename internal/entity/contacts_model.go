package entity

// GlobalContact сущность контакта
type GlobalContact struct {
	ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID int    `gorm:"not null;index" json:"account_id"`
	Email     string `gorm:"size:255;not null" json:"email"`
	Name      string `gorm:"size:255" json:"name"`
	Status    string `gorm:"size:50" json:"status"`
}
