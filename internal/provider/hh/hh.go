package hh

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/lookandhqte/workHelper/config"
	"github.com/lookandhqte/workHelper/internal/entity"
)

type Provider struct {
	app    *entity.App
	client *http.Client
}

// New создает нового провайдера
func New() *Provider {
	cfg := config.Load()
	return &Provider{
		app: &entity.App{
			RedirectURI:  cfg.RedirectURI,
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
		},
		client: &http.Client{},
	}
}

// GetToken меняет auth code на пару токенов, возвращает токены
func (r *Provider) GetToken(code string) (*entity.Token, error) {
	data := &url.Values{}
	data.Set("client_id", r.app.ClientID)
	data.Set("client_secret", r.app.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", r.app.RedirectURI)

	req, err := http.NewRequest(http.MethodPost, "https://api.hh.ru/token", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Printf("err while rnew request func get token hh.go: %v\n", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := r.client.Do(req)
	if err != nil {
		fmt.Printf("err while do req func get tokens hh.go: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err while readall get token func %v\n", err)
		return nil, err
	}

	responseData := &entity.Token{}
	if err := json.Unmarshal(body, responseData); err != nil {
		fmt.Printf("err while unmarshal body: %v\n", err)
		return nil, err
	}

	if responseData.AccessToken == "" && responseData.ExpiresIn == 0 && responseData.RefreshToken == "" {
		fmt.Printf("null result response data: %v\n", err)
		return nil, err
	}

	return responseData, nil
}
