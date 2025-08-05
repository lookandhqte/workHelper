package persistent

import (
	"gorm.io/gorm"
)

type DatabaseStorage struct {
	DB *gorm.DB
}

func NewDatabaseStorage() (*DatabaseStorage, error) {
	// db, err := gorm.Open().Open()
	return &DatabaseStorage{}, nil
}
