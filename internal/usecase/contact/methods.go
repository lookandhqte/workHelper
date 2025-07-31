package contact

import "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"

func (uc *ContactUseCase) GetTokens() (*entity.Token, error) {
	return uc.repo.GetTokens()
}
