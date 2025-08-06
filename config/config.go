package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

//Config структура конфига
type Config struct {
	JWTSecret   string
	HTTPAddr    string
	RedirectURI string
	DSN         string
	StorageType string
}

//Load подгрузка .env и создание конфига
func Load() *Config {
	if err := godotenv.Load("./.env"); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		JWTSecret:   getEnv("JWT_SECRET", ""),
		HTTPAddr:    getEnv("PORT", ":2020"),
		RedirectURI: getEnv("REDIRECT_URI", ""),
		DSN:         getEnv("DSN", ""),
		StorageType: getEnv("STORAGE_TYPE", "in-memory"),
	}
}

//getEnv получение переменных окружения
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		fmt.Println(value)
		return value
	}
	return defaultValue
}
