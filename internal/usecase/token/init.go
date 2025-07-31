package token

import "amocrm_golang/internal/entity"

type TokenUseCase struct {
	repo tokenRepo
}

type tokenRepo interface {
	AddTokens(response *entity.Token) error
	GetRefreshToken() (string, error)
	DeleteTokens() error
	UpdateTokens(response *entity.Token) error
}

func New(r tokenRepo) *TokenUseCase {
	return &TokenUseCase{repo: r}
}
