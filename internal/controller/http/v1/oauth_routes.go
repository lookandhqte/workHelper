package v1

import (
	"amocrm_golang/config"
	"amocrm_golang/internal/entity"
	"amocrm_golang/internal/usecase/token"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type oauthRoutes struct {
	uc token.TokenUseCase
}

type errResponse struct {
	Error string `json:"error"`
}

func NewTokenRoutes(handler *gin.RouterGroup, uc token.TokenUseCase) {
	r := &oauthRoutes{uc}

	h := handler.Group("/oauth")
	{
		h.GET("/", r.getTokens)
		h.PUT("/update", r.updateToken)
		h.DELETE("/", r.deleteTokens)
		h.GET("/redirect", r.handleRedirect)
	}

	r.StartTokenRefresher()
}

func (r *oauthRoutes) StartTokenRefresher() {
	ticker := time.NewTicker(1 * time.Hour) // Проверяем каждый час
	go func() {
		for range ticker.C {
			tokens, err := r.uc.GetTokens()
			if err != nil || tokens == nil {
				continue
			}

			// Обновляем токен, если до истечения осталось меньше 1 часа
			expiryTime := tokens.ServerTime + tokens.ExpiresIn

			if int(time.Now().Unix())-expiryTime <= 3600 {
				cfg := config.Load()
				newTokens, err := r.UpdateTokens(cfg.ClientID, cfg)
				if err != nil {
					log.Printf("Failed to refresh token: %v", err)
					continue
				}

				if err := r.uc.UpdateTokens(newTokens); err != nil {
					log.Printf("Failed to save refreshed tokens: %v", err)
				}
			}
		}
	}()
}

func (r *oauthRoutes) updateToken(c *gin.Context) {
	cfg := config.Load()
	id := c.Query("client_id")
	resp, err := r.UpdateTokens(id, cfg)

	if err != nil {
		err_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := r.uc.UpdateTokens(resp); err != nil {
		err_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *oauthRoutes) UpdateTokens(client_id string, config *config.Config) (*entity.Token, error) {
	refresh, err := r.uc.GetRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("error while getting refresh token: %v", err)
	}

	data := url.Values{}
	data.Set("client_id", client_id)
	data.Set("client_secret", config.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refresh)
	data.Set("redirect_uri", config.RedirectURI)

	return r.sendTokenRequest(data, config.AccessTokenURL)
}

func (r *oauthRoutes) handleRedirect(c *gin.Context) {
	// Извлекаем параметры из запроса
	code := c.Query("code")
	clientID := c.Query("client_id")
	if code == "" {
		c.JSON(http.StatusBadRequest, errResponse{Error: "Authorization code is required"})
		return
	}

	cfg := config.Load()

	tokens, err := r.GetTokensByAuthCode(code, clientID, cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse{Error: err.Error()})
		return
	}

	if err := r.uc.Create(tokens); err != nil {
		c.JSON(http.StatusInternalServerError, errResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"tokens": tokens,
	})
}

func (r *oauthRoutes) GetTokensByAuthCode(code string, client_id string, config *config.Config) (*entity.Token, error) {
	data := url.Values{}
	data.Set("client_id", client_id)
	data.Set("client_secret", config.ClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", config.RedirectURI)

	return r.sendTokenRequest(data, config.AccessTokenURL)
}

func (r *oauthRoutes) getTokens(c *gin.Context) {
	cfg := config.Load()
	tokens, err := r.GetTokensByAuthCode(c.Query("code"), c.Query("client_id"), cfg)

	if err != nil {
		err_Response(c, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Printf("Tokens:\n%s\n%s", tokens.AccessToken, tokens.RefreshToken)

	if err := r.uc.Create(tokens); err != nil {
		err_Response(c, http.StatusInternalServerError, err.Error())
		return
	}
}

func (r *oauthRoutes) deleteTokens(c *gin.Context) {

	if err := r.uc.Delete(); err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func err_Response(c *gin.Context, code int, err string) {
	c.AbortWithStatusJSON(code, errResponse{Error: err})
}

func prepareRequest(url string, data url.Values) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func sendRequest(req *http.Request) ([]byte, error) {
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

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func parseTokenResponse(body []byte) (*entity.Token, error) {
	responseData := &entity.Token{}
	if err := json.Unmarshal(body, responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

func (r *oauthRoutes) sendTokenRequest(data url.Values, url string) (*entity.Token, error) {
	req, err := prepareRequest(url, data)
	if err != nil {
		return nil, fmt.Errorf("request preparation failed: %v", err)
	}

	body, err := sendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	token, err := parseTokenResponse(body)
	if err != nil {
		return nil, fmt.Errorf("response parsing failed: %v", err)
	}

	return token, nil
}
