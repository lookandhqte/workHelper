package token

import "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"

func (uc *TokenUseCase) Create(response *entity.Token) error {
	return uc.repo.AddTokens(response)
}
func (uc *TokenUseCase) UpdateTokens(resp *entity.Token) error {
	return uc.repo.UpdateTokens(resp)
}

func (uc *TokenUseCase) GetRefreshToken() (string, error) {
	return uc.repo.GetRefreshToken()
}

func (uc *TokenUseCase) GetTokens() (*entity.Token, error) {
	return uc.repo.GetTokens()
}

func (uc *TokenUseCase) GetConst(req string) (int, error) {
	return uc.repo.GetConst(req)
}

func (uc *TokenUseCase) Delete() error {
	return uc.repo.DeleteTokens()
}
