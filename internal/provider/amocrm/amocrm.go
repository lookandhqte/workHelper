package amocrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

const (
	baseURL        = "https://spetser.amocrm.ru/"
	SlicesCapacity = 10
)

type Provider struct {
	client *http.Client
}

// New создает нового провайдера
func New() *Provider {
	return &Provider{
		client: &http.Client{},
	}
}

//GetTokensByAuuthCode получает токены с кодом auth
func (r *Provider) GetTokens(integration *entity.Integration) error {

	data := r.PrepareData("authorization_code", integration)

	fullURL := r.MakeRouteURL("/oauth2/access_token")
	req, err := r.PreparePostRequest(fullURL, data)
	if err != nil {
		return err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("API error: %d, body: %s", resp.StatusCode, string(body))
	}

	responseData := &entity.Token{}
	if err := json.Unmarshal(body, responseData); err != nil {
		return err
	}

	responseData.IntegrationID = integration.ID

	integration.Token = responseData

	return nil
}

//GetContacts возвращает контакты
func (r *Provider) GetContacts(token string) (*[]entity.Contact, error) {
	fullURL := r.MakeRouteURL("/api/v4/contacts")

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	body, err := r.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while sending request to get contacts")
	}

	var apiResponse APIContactsResponseDTO

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return apiResponse.ProcessContacts(), nil

}

//UpdateTokens обновляет токены
func (r *Provider) UpdateTokens(integration *entity.Integration) error {
	data := r.PrepareData("refresh_token", integration)

	fullURL := r.MakeRouteURL("/oauth2/access_token")
	token, err := r.SendTokenRequest(data, fullURL)
	if err != nil {
		log.Printf("error func update tokens -> send token req: %v", err)
		return err
	}
	integration.Token = token
	return nil
}

//ParseTokenResponse парсит ответ в токены
func (r *Provider) ParseTokenResponse(body []byte) (*entity.Token, error) {
	responseData := &entity.Token{}
	if err := json.Unmarshal(body, responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

//SendTokenRequest отправляет запрос на получение токенов
func (r *Provider) SendTokenRequest(data url.Values, url string) (*entity.Token, error) {
	req, err := r.PreparePostRequest(url, data)
	if err != nil {
		return nil, fmt.Errorf("request preparation failed: %v", err)
	}

	body, err := r.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	token, err := r.ParseTokenResponse(body)
	if err != nil {
		return nil, fmt.Errorf("response parsing failed: %v", err)
	}

	return token, nil
}

//PrepareData готовит url.Values для запроса
func (r *Provider) PrepareData(datacase string, integration *entity.Integration) url.Values {
	data := url.Values{}
	switch datacase {
	case "authorization_code":
		data.Set("client_id", integration.ClientID)
		data.Set("client_secret", integration.SecretKey)
		data.Set("grant_type", "authorization_code")
		data.Set("code", integration.AuthCode)
		data.Set("redirect_uri", integration.RedirectURL)
	case "refresh_token":
		data.Set("client_id", integration.ClientID)
		data.Set("client_secret", integration.SecretKey)
		data.Set("grant_type", "refresh_token")
		data.Set("refresh_token", integration.Token.RefreshToken)
		data.Set("redirect_uri", integration.RedirectURL)
	}
	return data
}

//MakeRouteURL возвращает полный URL адрес
func (r *Provider) MakeRouteURL(pathi string) string {

	base, _ := url.Parse(baseURL)
	base.Path = path.Join(base.Path, pathi)
	fullURL := base.String()

	return fullURL

}

//PreparePostRequest готовит post запрос
func (r *Provider) PreparePostRequest(url string, data url.Values) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

//SendRequest отправляет запрос
func (r *Provider) SendRequest(req *http.Request) ([]byte, error) {
	resp, err := r.client.Do(req)
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
