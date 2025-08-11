package entity

//Account структура аккаунта
type Account struct {
	ID              int           `json:"id" gorm:"primaryKey"`
	CacheExpires    int           `json:"cache_expires"`
	CreatedAt       int           `json:"created_at"`
	Integrations    []Integration `json:"integrations" gorm:"constraint:OnDelete:CASCADE;"`
	AccountContacts []Contact     `json:"contacts" gorm:"constraint:OnDelete:CASCADE;"`
}
