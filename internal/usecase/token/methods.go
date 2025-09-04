package token

import entity "github.com/lookandhqte/workHelper/internal/entity"

// Create создает токен
func (uc *UseCase) Create(token *entity.Token) error {
	return uc.repo.AddToken(token)
}

// ReturnExpiry возвращает expires_in токена
func (uc *UseCase) ReturnExpiry() (int, error) {
	return uc.repo.GetTokenExpiry()
}
