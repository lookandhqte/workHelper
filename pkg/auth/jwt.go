package auth

import (
	"time"

	config "git.amocrm.ru/gelzhuravleva/amocrm_golang/config"

	"github.com/golang-jwt/jwt/v4"
)

const (
	AccessTokenExpiry  = 86400 * time.Second
	RefreshTokenExpiry = 2592000 * time.Second
	SecretKey          = "amocrm_meow"
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

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
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
