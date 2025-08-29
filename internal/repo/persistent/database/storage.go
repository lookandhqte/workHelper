package database

import (
	"github.com/lookandhqte/workHelper/config"
	"github.com/lookandhqte/workHelper/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Storage структура
type Storage struct {
	DB *gorm.DB
}

// NewDatabaseStorage создает новое хранилище (база данных)
func NewDatabaseStorage(cfg *config.Config) (*Storage, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&entity.Account{},
		&entity.Token{},
	)
	if err != nil {
		return nil, err
	}

	return &Storage{DB: db}, nil
}
