package integration

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

//UseCase структура
type UseCase struct {
	repo integrationRepo
}

//integrationRepo абстракция для определения методов репозитория
type integrationRepo interface {
	AddIntegration(integration *entity.Integration) error
	GetIntegrations() (*[]entity.Integration, error)
	GetIntegration(id int) (*entity.Integration, error)
	UpdateIntegration(integration *entity.Integration) error
	DeleteIntegration(accountID int) error
	ReturnByClientID(clientID string) (*entity.Integration, error)
	UpdateToken(token *entity.Token) error
	GetTokens(id int) (*entity.Token, error)
}

//New создает новый репозиторий IntegrationUseCase
func New(r integrationRepo) *UseCase {
	return &UseCase{repo: r}
}
