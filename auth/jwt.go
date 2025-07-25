package auth

import (
	"time"

	config "amocrm_golang/confg"

	"github.com/golang-jwt/jwt/v4"
)

// Константы времени жизни токенов
const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
	SecretKey          = "amocrm_meow" // Дефолтный секретный ключ (переопределяется из .env)
)

var cfg = config.Load()

// Структура для хранения данных в JWT токене
type Claims struct {
	AccountID int `json:"account_id"`
	jwt.RegisteredClaims
}

// Создает новый JWT токен
func GenerateJWT(accountID int, expiry time.Duration) (string, error) {

	claims := &Claims{
		AccountID: accountID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// ParseJWT разбирает и валидирует JWT токен
func ParseJWT(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil // Используем секрет из конфига
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
