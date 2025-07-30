package token

import "amocrm_golang/internal/entity"

func (uc *TokenUseCase) Create(response *entity.Token) error {
	return uc.repo.AddTokens(response)
}

func (uc *TokenUseCase) UpdateRefreshToken(token string) error {
	return uc.repo.UpdateRToken(token)
}
func (uc *TokenUseCase) GetRefreshToken() (string, error) {
	return uc.repo.GetRefreshToken()
}
func (uc *TokenUseCase) UpdateAccessToken(token string) error {
	return uc.repo.UpdateAToken(token)
}

func (uc *TokenUseCase) Delete() error {
	return uc.repo.DeleteTokens()
}
