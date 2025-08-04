package persistent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/pkg/cache"
)

type MemoryStorage struct {
	mu               sync.RWMutex
	accounts         map[int]*entity.Account
	integrations     map[int]*entity.Integration
	main_account     *entity.Account
	main_integration *entity.Integration
	lastAccountID    int
	cache            *cache.Cache
}

const (
	BASE_ID_FOR_TOKENS  = 0
	ACCESS_EXPIRES_SEC  = 86400
	REFRESH_EXPIRES_SEC = 2592000
	CACHE_EXPIRES_SEC   = 604800
	BASE_URL            = "https://spetser.amocrm.ru/"
)

func NewMemoryStorage(c *cache.Cache) *MemoryStorage {
	return &MemoryStorage{
		accounts:      make(map[int]*entity.Account),
		integrations:  make(map[int]*entity.Integration),
		lastAccountID: 0,
		cache:         c,
	}
}

func (m *MemoryStorage) AddAccount(account *entity.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastAccountID++
	account.ID = m.lastAccountID
	account.CacheExpires = account.CreatedAt + CACHE_EXPIRES_SEC
	account.RefreshTokenExpiresIn = account.CreatedAt + REFRESH_EXPIRES_SEC
	account.AccessTokenExpiresIn = account.CreatedAt + ACCESS_EXPIRES_SEC

	m.accounts[account.ID] = account

	return nil
}

func (m *MemoryStorage) GetAccounts() ([]*entity.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	accounts := make([]*entity.Account, 0, len(m.accounts))
	for _, account := range m.accounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (m *MemoryStorage) GetAccount(id int) (*entity.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	account, exists := m.accounts[id]
	if !exists {
		return nil, fmt.Errorf("account not found to get account")
	}

	return account, nil
}

func (m *MemoryStorage) GetMainAccount() *entity.Account {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.main_account
}

func (m *MemoryStorage) GetMainIntegration() *entity.Integration {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.main_integration
}

func (m *MemoryStorage) ChangeMainAccount(new_id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	account, err := m.GetAccount(new_id)
	if err != nil {
		return err
	}
	m.main_account = account
	return nil
}

func (m *MemoryStorage) ChangeMainIntegration(new_id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	integration, err := m.GetIntegration(new_id)
	if err != nil {
		return err
	}
	m.main_integration = integration
	return nil
}

func (m *MemoryStorage) GetIntegration(id int) (*entity.Integration, error) {
	m.mu.Lock()
	m.mu.Unlock()
	integration, exists := m.integrations[id]
	if !exists {
		return nil, fmt.Errorf("no integrations with these id")
	}

	return integration, nil

}

func (m *MemoryStorage) GetMainAccountTokens() *entity.Token {
	m.mu.Lock()
	defer m.mu.Unlock()

	var tokens *entity.Token
	tokens.AccessToken = m.main_account.AccessToken
	tokens.RefreshToken = m.main_account.RefreshToken
	tokens.ExpiresIn = m.main_account.RefreshTokenExpiresIn

	return tokens
}

func (m *MemoryStorage) GetAccountWithCache(id int) (*entity.Account, error) {
	if cached, ok := m.cache.Get(id); ok {
		return cached.(*entity.Account), nil
	}

	account, err := m.GetAccount(id)
	if err != nil {
		return nil, err
	}

	m.cache.Set(id, account, time.Hour)
	return account, nil
}

func (m *MemoryStorage) GetTokensByAuthCode(code string, client_id string) (*entity.Token, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	integration, err := m.GetIntegrationsByClientID(client_id)
	if err != nil {
		return nil, fmt.Errorf("error n func get integr by client id -> error in get tokens method")
	}
	data := url.Values{}
	data.Set("client_id", client_id)
	data.Set("client_secret", integration.SecretKey)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", integration.RedirectUrl)
	base, err := url.Parse(BASE_URL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}
	base.Path = path.Join(base.Path, "/oauth2/access_token")
	fullURL := base.String()

	req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("API error: %d, body: %s", resp.StatusCode, string(body))
	}

	responseData := &entity.Token{}
	if err := json.Unmarshal(body, responseData); err != nil {
		return nil, err
	}

	return responseData, nil
}

func (m *MemoryStorage) UpdateAccount(account *entity.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.accounts[account.ID]; !exists {
		return fmt.Errorf("account not found to update")
	}

	m.accounts[account.ID] = account
	return nil
}

func (m *MemoryStorage) UpdateAccountTokens(tokens *entity.Token) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.main_account.AccessToken = tokens.AccessToken
	m.main_account.AccessTokenExpiresIn = int(time.Now().Unix()) + ACCESS_EXPIRES_SEC
	m.main_account.RefreshToken = tokens.RefreshToken

	m.accounts[m.main_account.ID] = m.main_account
	return nil
}

func (m *MemoryStorage) DeleteAccount(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.accounts[id]; !exists {
		return fmt.Errorf("account not found to delete")
	}

	delete(m.accounts, id)
	return nil
}

func (m *MemoryStorage) AddIntegration(integration *entity.Integration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.integrations[integration.AccountID] = integration
	return nil
}

func (m *MemoryStorage) GetIntegrations() ([]*entity.Integration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	integrations := make([]*entity.Integration, 0, len(m.integrations))
	for _, integration := range m.integrations {
		integrations = append(integrations, integration)
	}

	return integrations, nil
}

func (m *MemoryStorage) GetIntegrationsByClientID(client_id string) (*entity.Integration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var integrationEntity *entity.Integration
	for _, integration := range m.integrations {
		if integration.ClientID == client_id {
			integrationEntity = integration
		}
	}

	if integrationEntity.AuthCode == "" {
		return nil, fmt.Errorf("no such client_id to get intgrations by client_id")
	}

	return integrationEntity, nil
}

func (m *MemoryStorage) GetAccountIntegrations(accountID int) (*entity.Integration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	integration, exists := m.integrations[accountID]
	if !exists {
		return nil, fmt.Errorf("integration not found to get account integrations")
	}

	return integration, nil
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

func (m *MemoryStorage) AddTokens(response *entity.Token) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache.SetToken(BASE_ID_FOR_TOKENS, response, time.Duration(ACCESS_EXPIRES_SEC)*time.Second)

	return nil
}

func (m *MemoryStorage) DeleteTokens() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache.DeleteToken(BASE_ID_FOR_TOKENS)

	return nil
}

func (m *MemoryStorage) UpdateTokens(response *entity.Token) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache.Set(BASE_ID_FOR_TOKENS, response, time.Duration(ACCESS_EXPIRES_SEC)*time.Second)
	return nil
}

func (m *MemoryStorage) GetRefreshToken() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ref, exists := m.cache.GetToken(BASE_ID_FOR_TOKENS)
	if !exists {
		return "", fmt.Errorf("no refresh key in storage")
	}
	return ref.RefreshToken, nil
}

func (m *MemoryStorage) GetTokens() (*entity.Token, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Println(m)
	val, exists := m.cache.GetToken(BASE_ID_FOR_TOKENS)
	if !exists {
		return nil, fmt.Errorf("no tokens in storage")
	}
	return val, nil
}
