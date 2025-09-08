package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config структура конфига
type Config struct {
	HTTPAddr      string
	RedirectURI   string
	ClientID      string
	ClientSecret  string
	DSN           string
	StorageType   string
	WorkerAmount  string
	BeanstalkAddr string
	WorkerDSN     string
	DeepseekAPI   string
}

// Load подгрузка .env и создание конфига
func Load() *Config {
	if err := godotenv.Load("./.env"); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		HTTPAddr:      getEnv("PORT", ":2020"),
		RedirectURI:   getEnv("REDIRECT_URI", ""),
		DSN:           getEnv("DSN", ""),
		StorageType:   getEnv("STORAGE_TYPE", "in-memory"),
		ClientID:      getEnv("CLIENT_ID", ""),
		ClientSecret:  getEnv("CLIENT_SECRET", ""),
		WorkerAmount:  getEnv("WORKER_AMOUNT", ""),
		BeanstalkAddr: getEnv("BEANSTALK_ADDR", ""),
		WorkerDSN:     getEnv("WORKER_DSN", ""),
		DeepseekAPI:   getEnv("DEEPSEEK_API", ""),
	}
}

// getEnv получение переменных окружения
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
