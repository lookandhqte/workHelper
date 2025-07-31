package token

import "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"

type TokenUseCase struct {
	repo tokenRepo
}

type tokenRepo interface {
	AddTokens(response *entity.Token) error
	GetTokens() (*entity.Token, error)
	GetRefreshToken() (string, error)
	DeleteTokens() error
	UpdateTokens(response *entity.Token) error
	GetConst(req string) (int, error)
}

func New(r tokenRepo) *TokenUseCase {
	return &TokenUseCase{repo: r}
}
