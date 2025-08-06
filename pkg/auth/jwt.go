package auth

import (
	"time"

	config "git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	"github.com/golang-jwt/jwt/v4"
)

const (
	//AccessTokenExpiry определяет через сколько обновлять access токен
	AccessTokenExpiry = 86400 * time.Second
	//RefreshTokenExpiry определяет через сколько обновлять access токен
	RefreshTokenExpiry = 2592000 * time.Second
	//SecretKey определяет секретный ключ
	SecretKey = "amocrm_meow"
)

//var cfg переменная конфигурации
var cfg = config.Load()

//Claims структура
type Claims struct {
	AccountID int `json:"account_id"`
	jwt.RegisteredClaims
}

//GenerateJWT генерирует токены
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

//ParseJWT парсит токены
func ParseJWT(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(_ *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
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
