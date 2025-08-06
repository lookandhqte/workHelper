package storage

import (
	"log"

	config "git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	database "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/repo/persistent/database"
	inmemory "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/repo/persistent/inmemory"
	cache "git.amocrm.ru/gelzhuravleva/amocrm_golang/pkg/cache"
)

//NewStorage создает новое хранилище в зависимости от STORAGE_TYPE
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
