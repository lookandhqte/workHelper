package inmemory

import (
	"fmt"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

func (m *MemoryStorage) AddIntegration(integration *entity.Integration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.integrations[integration.AccountID] = integration

	m.accounts[integration.AccountID].Integrations = append(m.accounts[integration.AccountID].Integrations, *integration)
	return nil
}

func (m *MemoryStorage) GetIntegration(id int) (*entity.Integration, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	integration, exists := m.integrations[id]

	if !exists {
		return nil, fmt.Errorf("no integrations with these id")
	}

	return integration, nil

}

func (m *MemoryStorage) GetIntegrations() (*[]entity.Integration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	integrations := make([]entity.Integration, 0, len(m.integrations))
	for _, integration := range m.integrations {
		integrations = append(integrations, *integration)
	}

	return &integrations, nil
}

func (m *MemoryStorage) UpdateIntegration(integration *entity.Integration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.integrations[integration.AccountID]; !exists {
		return fmt.Errorf("integration not found to update")
	}

	m.integrations[integration.AccountID] = integration
	return nil
}

func (m *MemoryStorage) DeleteIntegration(accountID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.integrations[accountID]; !exists {
		return fmt.Errorf("integration not found to delete integration")
	}

	delete(m.integrations, accountID)
	return nil
}

func (m *MemoryStorage) ReturnByClientID(client_id string) (*entity.Integration, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	integrations, _ := m.GetIntegrations()
	for _, integration := range *integrations {
		if integration.ClientID == client_id {
			return &integration, nil
		}
	}

	return nil, fmt.Errorf("haven't found here your integration")
}
