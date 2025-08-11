package unisender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

const (
	baseURL = "https://api.unisender.com"
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

//getListsFromUnisender
func (r *Provider) GetLists(integrations *[]entity.Integration) error {

	integrationsPtr := *integrations

	var unisenderKey string
	for _, integration := range integrationsPtr {
		if integration.Token.UnisenderKey != "" {
			unisenderKey = integration.Token.UnisenderKey
			break
		}
	}

	fullURL := baseURL + "/ru/api/getLists?format=json&api_key=" + unisenderKey

	var data url.Values = url.Values{}
	req, err := http.NewRequest(http.MethodGet, fullURL, bytes.NewBufferString(data.Encode()))

	if err != nil {
		log.Printf("error while new request -> get lists func")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	body, err := r.SendRequest(req)
	if err != nil {
		log.Printf("error while sending request func auth unisender: %v", err)
	}
	defer req.Body.Close()
	responseData := &ListUnisenderDTO{}

	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Printf("error while unmarshal data func auth unisender: %v", err)
	}

	return nil
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
