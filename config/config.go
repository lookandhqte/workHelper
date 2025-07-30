package config

import (
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
	ClientID       string
	ClientSecret   string
	RedirectURI    string
	AccessTokenURL string
}

//Подгрузка .env и создание конфига
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		JWTSecret:      getEnv("JWT_SECRET", ""),
		AccessTokenTTL: time.Minute * 15,
		HTTPAddr:       getEnv("PORT", ":2020"),
		ClientID:       getEnv("CLIENT_ID", ""),
		ClientSecret:   getEnv("CLIENT_SECRET", ""),
		RedirectURI:    getEnv("REDIRECT_URI", ""),
		AccessTokenURL: getEnv("ACCESS_TOKEN_URL", ""),
	}
}

//Получение переменных окружения .env
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
