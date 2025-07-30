package token

import "amocrm_golang/internal/entity"

type TokenUseCase struct {
	repo tokenRepo
}

type tokenRepo interface {
	AddTokens(response *entity.Token) error
	GetRefreshToken() (string, error)
	DeleteTokens() error
	UpdateRToken(refresh string) error
	UpdateAToken(access string) error
}

func New(r tokenRepo) *TokenUseCase {
	return &TokenUseCase{repo: r}
}
