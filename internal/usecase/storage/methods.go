package storage

import (
	"log"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/repo/persistent/database"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/repo/persistent/inmemory"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/pkg/cache"
)

func NewStorage(c *cache.Cache, cfg *config.Config) *Storage {
	switch cfg.StorageType {
	case "in-memory":
		DB = inmemory.NewMemoryStorage(c)
		return &DB
	case "database":
		db, err := database.NewDatabaseStorage(cfg)
		if err != nil {
			log.Printf("error in new database storage func -> error in new storage func: %v", err)
			DB = inmemory.NewMemoryStorage(c)
			log.Printf("using in-memory cause of err")
			return &DB
		}
		log.Printf("using mysql database")
		DB = db
		return &DB
	default:
		DB = inmemory.NewMemoryStorage(c)
		log.Println("Using in-memory storage")
		return &DB
	}
}
