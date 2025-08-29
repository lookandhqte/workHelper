package storage

import (
	"log"

	config "github.com/lookandhqte/workHelper/config"
	database "github.com/lookandhqte/workHelper/internal/repo/persistent/database"
	inmemory "github.com/lookandhqte/workHelper/internal/repo/persistent/inmemory"
)

// NewStorage создает новое хранилище в зависимости от STORAGE_TYPE
func NewStorage(cfg *config.Config) *Storage {
	switch cfg.StorageType {
	case "in-memory":
		DB = inmemory.NewMemoryStorage()
		return &DB
	case "database":
		db, err := database.NewDatabaseStorage(cfg)
		if err != nil {
			log.Printf("error in new database storage func -> error in new storage func: %v", err)
			DB = inmemory.NewMemoryStorage()
			log.Printf("using in-memory cause of err")
			return &DB
		}
		log.Printf("using mysql database")
		DB = db
		return &DB
	default:
		DB = inmemory.NewMemoryStorage()
		log.Println("Using in-memory storage")
		return &DB
	}
}
