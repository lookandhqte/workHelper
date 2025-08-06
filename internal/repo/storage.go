package repo

import (
	"log"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/repo/persistent"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/pkg/cache"
)

type Storage interface {
	AddAccount(account *entity.Account) error
	GetAccounts() ([]*entity.Account, error)
	GetAccount(id int) (*entity.Account, error)
	GetAccountIntegrations(accountID int) (*[]entity.Integration, error)
	UpdateAccount(account *entity.Account) error
	DeleteAccount(id int) error
	AddIntegration(integration *entity.Integration) error
	GetIntegration(id int) (*entity.Integration, error)
	GetIntegrations() (*[]entity.Integration, error)
	UpdateIntegration(integration *entity.Integration) error
	DeleteIntegration(accountID int) error
	ReturnByClientID(client_id string) (int, error)
}

var DB Storage

func NewStorage(c *cache.Cache, cfg *config.Config) *Storage {
	switch cfg.StorageType {
	case "in-memory":
		DB = persistent.NewMemoryStorage(c)
		return &DB
	case "database":
		db, err := persistent.NewDatabaseStorage(cfg)
		if err != nil {
			log.Printf("error in new database storage func -> error in new storage func: %v", err)
			DB = persistent.NewMemoryStorage(c)
			log.Printf("using in-memory cause of err")
			return &DB
		}
		log.Printf("using mysql database")
		DB = db
		return &DB
	default:
		DB = persistent.NewMemoryStorage(c)
		log.Println("Using in-memory storage")
		return &DB
	}
}
