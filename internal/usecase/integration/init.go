package integration

import (
	"context"
	"sync"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/dto"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

type IntegrationUseCase struct {
	repo integrationRepo
}

type integrationRepo interface {
	AddIntegration(integration *entity.Integration) error
	GetIntegrations() ([]*entity.Integration, error)
	GetIntegration(id int) (*entity.Integration, error)
	UpdateIntegration(integration *entity.Integration) error
	DeleteIntegration(accountID int) error
	GetIntegrationByClientID(client_id string) (*entity.Integration, error)
	GetContacts(token string) (*dto.ContactsResponse, error)
	Start(wg *sync.WaitGroup) func(ctx context.Context)
}

func New(r integrationRepo) *IntegrationUseCase {
	return &IntegrationUseCase{repo: r}
}
