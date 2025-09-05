package entity

// Account структура аккаунта
type Account struct {
	ID        int   `json:"id" gorm:"primaryKey"` //автоген
	CreatedAt int   `json:"created_at"`
	Token     Token `json:"account_tokens" gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE;"`
}
