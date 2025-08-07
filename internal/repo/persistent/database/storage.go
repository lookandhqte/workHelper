package database

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//Storage структура
type Storage struct {
	DB *gorm.DB
}

//NewDatabaseStorage создает новое хранилище (база данных)
func NewDatabaseStorage(cfg *config.Config) (*Storage, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	db.Migrator().DropTable(&entity.Token{})
	db.Migrator().DropTable(&entity.Contact{})
	db.Migrator().DropTable(&entity.Integration{})
	db.Migrator().DropTable(&entity.Account{})
	err = db.AutoMigrate(
		&entity.Account{},
		&entity.Integration{},
		&entity.Token{},
		&entity.Contact{},
	)
	if err != nil {
		return nil, err
	}

	return &Storage{DB: db}, nil
}
