package token

import (
	entity "github.com/lookandhqte/workHelper/internal/entity"
)

// UseCase структура
type UseCase struct {
	repo tokenRepo
}

// tokenRepo абстракция для определения методов репозитория
type tokenRepo interface {
	AddToken(token *entity.Token) error
	GetTokenExpiry() (int, error)
}

// New создает новый репозиторий
func New(r tokenRepo) *UseCase {
	return &UseCase{repo: r}
}
