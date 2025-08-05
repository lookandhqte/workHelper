package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

//Структура конфига
type Config struct {
	JWTSecret      string
	AccessTokenTTL time.Duration
	HTTPAddr       string
	RedirectURI    string
	DSN            string
	StorageType    string
}

//Подгрузка .env и создание конфига
func Load() *Config {
	if err := godotenv.Load("./.env"); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		JWTSecret:      getEnv("JWT_SECRET", ""),
		AccessTokenTTL: time.Minute * 15,
		HTTPAddr:       getEnv("PORT", ":2020"),
		RedirectURI:    getEnv("REDIRECT_URI", ""),
		DSN:            getEnv("DSN", ""),
		StorageType:    getEnv("STORAGE_TYPE", "in-memory"),
	}

}

//Получение переменных окружения .env
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		fmt.Println(value)
		return value
	}
	return defaultValue
}
