package auth

import (
	"time"

	config "amocrm_golang/confg"

	"github.com/golang-jwt/jwt/v4"
)

const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
	SecretKey          = "amocrm_meow" //добавить подгрузку из .env
)

var cfg = config.Load()

type Claims struct {
	AccountID int `json:"account_id"`
	jwt.RegisteredClaims
}

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

func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
