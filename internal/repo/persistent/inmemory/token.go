package inmemory

import "github.com/lookandhqte/workHelper/internal/entity"

// AddAccount создает аккаунт в in-memory хранилище
func (m *MemoryStorage) AddToken(token *entity.Token) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.token = token
	return nil
}

func (m *MemoryStorage) GetTokenExpiry() (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.token.ExpiresIn, nil
}
