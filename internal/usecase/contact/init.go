package contact

import "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"

type ContactUseCase struct {
	repo contactRepo
}

type contactRepo interface {
	GetTokens() (*entity.Token, error)
}

func New(r contactRepo) *ContactUseCase {
	return &ContactUseCase{repo: r}
}
