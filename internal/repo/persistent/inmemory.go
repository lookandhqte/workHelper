package persistent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/dto"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/pkg/cache"
)

type MemoryStorage struct {
	mu                  sync.RWMutex
	accounts            map[int]*entity.Account
	integrations        map[int]*entity.Integration
	active_account      *entity.Account
	active_integrations map[int]bool
	lastAccountID       int
	cache               *cache.Cache
}

const (
	BASE_ID_FOR_TOKENS  = 0
	ACCESS_EXPIRES_SEC  = 86400
	REFRESH_EXPIRES_SEC = 2592000
	CACHE_EXPIRES_SEC   = 604800
	BASE_URL            = "https://spetser.amocrm.ru/"
	REF_THRESHOLD_SEC   = 3600
)

func NewMemoryStorage(c *cache.Cache) *MemoryStorage {
	return &MemoryStorage{
		accounts:            make(map[int]*entity.Account),
		integrations:        make(map[int]*entity.Integration),
		active_account:      &entity.Account{},
		active_integrations: make(map[int]bool),
		lastAccountID:       0,
		cache:               c,
	}
}

//Доступные пользователю методы

//Методы аккаунта

func (m *MemoryStorage) AddAccount(account *entity.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastAccountID++
	account.ID = m.lastAccountID
	account.CacheExpires = account.CreatedAt + CACHE_EXPIRES_SEC

	m.accounts[account.ID] = account

	return nil
}

func (m *MemoryStorage) GetAccounts() ([]*entity.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	//Добавить возврат аккаунтов из кэша
	// accs, err := m.GetAccountWithCache()

	accounts := make([]*entity.Account, 0, len(m.accounts))
	for _, account := range m.accounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (m *MemoryStorage) GetAccount(id int) (*entity.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	acc, err := m.GetAccountWithCache(id)
	if err != nil {
		account, exists := m.accounts[id]
		if !exists {
			return nil, fmt.Errorf("account not found to get account")
		}
		return account, nil
	}

	return acc, nil
}

func (m *MemoryStorage) ChangeActiveAccount(new_id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if new_id == m.active_account.ID {
		return fmt.Errorf("this acc is active")
	}
	account, err := m.GetAccount(new_id)
	if err != nil {
		return err
	}
	m.active_account = account
	return nil
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

func (m *MemoryStorage) DeleteAccount(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.accounts[id]; !exists {
		return fmt.Errorf("account not found to delete")
	}

	delete(m.accounts, id)
	return nil
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

//Методы интеграций

func (m *MemoryStorage) GetIntegration(id int) (*entity.Integration, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	integration, exists := m.integrations[id]

	if !exists {
		return nil, fmt.Errorf("no integrations with these id")
	}

	return integration, nil

}

func (m *MemoryStorage) GetIntegrationByClientID(client_id string) (*entity.Integration, error) {
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

//Добавить появление интеграций в активном аккаунте. Не просто так же он есть....
func (m *MemoryStorage) AddIntegration(integration *entity.Integration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.integrations[integration.AccountID] = integration
	return nil
}

//Получение интеграций только из активного аккаунта
//Нужно еще сущность админа как-то сделать чтобы ему было видно все-все
func (m *MemoryStorage) GetIntegrations() ([]*entity.Integration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	integrations := make([]*entity.Integration, 0, len(m.integrations))
	for _, integration := range m.integrations {
		integrations = append(integrations, integration)
	}

	return integrations, nil
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

func (m *MemoryStorage) GetContacts(token string) (*dto.ContactsResponse, error) {
	fullURL := MakeRouteURL("/api/v4/contacts")

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	body, err := SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while sending request to get contacts")
	}

	var apiResponse dto.APIContactsResponse

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return apiResponse.ToContactsResponse(), nil

}

//Внутренние методы, не должны быть доступны пользователю

//Методы аккаунта

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

func (m *MemoryStorage) GetActiveAccount() *entity.Account {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.active_account
}

//Методы интеграции

func (m *MemoryStorage) UpdateTokens(client_id string) (*entity.Token, error) {
	integration, err := m.GetIntegrationByClientID(client_id)
	if err != nil {
		fmt.Print("error in func update tokens -> get int by client id")
	}

	data := url.Values{}
	data.Set("client_id", client_id)
	data.Set("client_secret", integration.SecretKey)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", integration.Token.RefreshToken)
	data.Set("redirect_uri", integration.RedirectUrl)
	fullUrl := MakeRouteURL(BASE_URL)
	return SendTokenRequest(data, fullUrl)
}

//Добавить логику выбора активной интеграции. ее айдишник будет фигурировать везде. Убрать метод гет инт по клиент id.

//Методы oauth

func (m *MemoryStorage) GetTokensByAuthCode(code string, client_id string) (*entity.Token, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	integration, err := m.GetIntegrationByClientID(client_id)
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

func (m *MemoryStorage) Start(wg *sync.WaitGroup) func(ctx context.Context) {
	return func(ctx context.Context) {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.refreshTokensBatch(wg)
			case <-ctx.Done():
				log.Println("Token refresher stopped")
				return
			}
		}
	}
}

func (m *MemoryStorage) refreshTokensBatch(wg *sync.WaitGroup) {
	integr, err := m.GetIntegrations()
	if err != nil {
		log.Printf("Failed to get active integrations: %v", err)
		return
	}

	sem := make(chan struct{}, 10)

	for i := range integr {
		wg.Add(1)
		sem <- struct{}{}

		go func(integration *entity.Integration) {
			defer wg.Done()
			defer func() { <-sem }()

			m.refreshTokenIfNeeded(integration)
		}(integr[i])
	}

	wg.Wait()
}

func (m *MemoryStorage) refreshTokenIfNeeded(integration *entity.Integration) {
	expiryTime := integration.Token.ServerTime + integration.Token.ExpiresIn
	now := time.Now().Unix()

	if expiryTime-int(now) <= REF_THRESHOLD_SEC {
		newTokens, err := m.UpdateTokens(integration.ClientID)
		if err != nil {
			log.Printf("[Acc:%d] Failed to refresh token: %v", integration.AccountID, err)
			return
		}
		//Нужно будет добавить шифрование токенов из auth.
		integration.Token = newTokens

		if err := m.UpdateIntegration(integration); err != nil {
			log.Printf("[Acc:%d] Failed to save tokens: %v", integration.AccountID, err)
		}
	}
}

//Методы запросов

func MakeRouteURL(pathi string) string {
	base, _ := url.Parse(BASE_URL)

	base.Path = path.Join(base.Path, pathi)
	fullURL := base.String()
	return fullURL
}

func PrepareRequest(url string, data url.Values) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func SendRequest(req *http.Request) ([]byte, error) {
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

	return body, nil
}

func ParseTokenResponse(body []byte) (*entity.Token, error) {
	responseData := &entity.Token{}
	if err := json.Unmarshal(body, responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

func SendTokenRequest(data url.Values, url string) (*entity.Token, error) {
	req, err := PrepareRequest(url, data)
	if err != nil {
		return nil, fmt.Errorf("request preparation failed: %v", err)
	}

	body, err := SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	token, err := ParseTokenResponse(body)
	if err != nil {
		return nil, fmt.Errorf("response parsing failed: %v", err)
	}

	return token, nil
}

// func (m *MemoryStorage) AddTokens(response *entity.Token) error {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	m.cache.SetToken(BASE_ID_FOR_TOKENS, response, time.Duration(ACCESS_EXPIRES_SEC)*time.Second)

// 	return nil
// }

// func (m *MemoryStorage) DeleteTokens() error {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	m.cache.DeleteToken(BASE_ID_FOR_TOKENS)

// 	return nil
// }

// func (m *MemoryStorage) UpdateTokens(response *entity.Token) error {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	m.cache.Set(BASE_ID_FOR_TOKENS, response, time.Duration(ACCESS_EXPIRES_SEC)*time.Second)
// 	return nil
// }

// func (m *MemoryStorage) GetTokens() (*entity.Token, error) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	fmt.Println(m)
// 	val, exists := m.cache.GetToken(BASE_ID_FOR_TOKENS)
// 	if !exists {
// 		return nil, fmt.Errorf("no tokens in storage")
// 	}
// 	return val, nil
// }
//Токены хранятся в интеграции, не в кэше...
