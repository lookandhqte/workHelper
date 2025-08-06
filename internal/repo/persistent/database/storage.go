package database

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//Storage структура
type Storage struct {
	DB *gorm.DB
}

//NewDatabaseStorage создает новое хранилище (база данных)
func NewDatabaseStorage(cfg *config.Config) (*Storage, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// db.Migrator().DropTable(&entity.Token{})
	// db.Migrator().DropTable(&entity.Integration{})
	// db.Migrator().DropTable(&entity.Account{})
	// db.Migrator().DropTable(&entity.Contact{})
	err = db.AutoMigrate(
		&entity.Account{},
		&entity.Integration{},
		&entity.Contact{},
		&entity.Token{},
	)
	if err != nil {
		return nil, err
	}

	return &Storage{DB: db}, nil
}
