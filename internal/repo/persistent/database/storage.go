package database

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseStorage struct {
	DB *gorm.DB
}

func NewDatabaseStorage(cfg *config.Config) (*DatabaseStorage, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(
		&entity.Account{},
		&entity.Integration{},
		&entity.Contact{},
		&entity.Token{},
	)
	if err != nil {
		return nil, err
	}

	return &DatabaseStorage{DB: db}, nil
}
