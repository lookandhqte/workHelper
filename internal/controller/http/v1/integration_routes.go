package v1

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
	"strconv"
	"sync"
	"time"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/dto"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/integration"

	"github.com/gin-gonic/gin"
)

type integrationRoutes struct {
	uc integration.IntegrationUseCase
}

const (
	REF_THRESHOLD_SEC = 3600
	BASE_Url          = "https://spetser.amocrm.ru/"
)

func NewIntegrationRoutes(handler *gin.RouterGroup, uc integration.IntegrationUseCase) {
	r := &integrationRoutes{uc}

	h := handler.Group("/integrations")
	{
		h.POST("/", r.createIntegration)
		h.GET("/", r.getIntegrations)
		h.PUT("/:id", r.updateIntegration)
		h.DELETE("/:id", r.deleteIntegration)
		h.GET("/redirect", r.handleRedirect)
		h.GET("/contacts", r.getContacts)
	}
	//r.StartTokenRefresher(context.Background())
}

func (r *integrationRoutes) getContacts(c *gin.Context) {
	integr, err := r.uc.GetIntegrationByClientID(c.Query("client_id"))

	if err != nil {
		error_Response(c, http.StatusUnauthorized, "error in  get contacts func -> get int by client id")
		return
	}

	if integr.Token.AccessToken == "" {
		error_Response(c, http.StatusUnauthorized, "access token missing")
		return
	}

	contacts, err := r.GetContacts(integr.Token.AccessToken)
	if err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"contacts": contacts,
	})

}

func (r *integrationRoutes) GetContacts(token string) (*dto.ContactsResponse, error) {
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

func (r *integrationRoutes) StartTokenRefresher(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Token refresher stopped")
			return
		case <-ticker.C:
			r.refreshTokensBatch()
		}
	}
}

func (r *integrationRoutes) refreshTokensBatch() {
	integr, err := r.uc.GetActiveIntegrations()
	if err != nil {
		log.Printf("Failed to get active integrations: %v", err)
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	for i := range integr {
		wg.Add(1)
		sem <- struct{}{}
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic in token refresh: %v", r)
			}
		}()
		go func(integration *entity.Integration) {
			defer wg.Done()
			defer func() { <-sem }()

			r.refreshTokenIfNeeded(integration)
		}(integr[i])
	}

	wg.Wait()
}

func (r *integrationRoutes) refreshTokenIfNeeded(integration *entity.Integration) {
	expiryTime := integration.Token.ServerTime + integration.Token.ExpiresIn
	now := time.Now().Unix()

	if expiryTime-int(now) <= REF_THRESHOLD_SEC {
		newTokens, err := r.UpdateTokens(integration.ClientID)
		if err != nil {
			log.Printf("[Acc:%s] Failed to refresh token: %v", integration.AccountID, err)
			return
		}

		if err := r.uc.UpdateTokens(newTokens); err != nil {
			log.Printf("[Acc:%s] Failed to save tokens: %v", integration.AccountID, err)
		}
	}
}

func (r *integrationRoutes) UpdateTokens(client_id string) (*entity.Token, error) {
	integration, err := r.uc.GetIntegrationByClientID(client_id)
	if err != nil {
		fmt.Print("error in func update tokens -> get int by client id")
	}

	data := url.Values{}
	data.Set("client_id", client_id)
	data.Set("client_secret", integration.SecretKey)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", integration.Token.RefreshToken)
	data.Set("redirect_uri", integration.RedirectUrl)
	fullUrl := MakeRouteURL(BASE_Url)
	return SendTokenRequest(data, fullUrl)
}

func (r *integrationRoutes) createIntegration(c *gin.Context) {
	var integration entity.Integration
	if err := c.ShouldBindJSON(&integration); err != nil {
		error_Response(c, http.StatusBadRequest, err.Error())
		return
	}
	//нужно присваивать интеграцию активному аккаунту и все действия с интеграциями от имени активного аккаунта
	if integration.AccountID == 0 {
		error_Response(c, http.StatusBadRequest, "account ID is required")
		return
	}

	if err := r.uc.Create(&integration); err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, integration)
}

func (r *integrationRoutes) getIntegrations(c *gin.Context) {
	integrations, err := r.uc.Return(nil)
	if err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, integrations)
}

func (r *integrationRoutes) updateIntegration(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_Response(c, http.StatusBadRequest, "ID must be integer")
		return
	}

	var integration entity.Integration
	if err := c.ShouldBindJSON(&integration); err != nil {
		error_Response(c, http.StatusBadRequest, err.Error())
		return
	}

	integration.AccountID = id
	if err := r.uc.Update(&integration); err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, integration)
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

func MakeRouteURL(pathi string) string {
	base, _ := url.Parse(BASE_Url)

	base.Path = path.Join(base.Path, pathi)
	fullURL := base.String()
	return fullURL
}

func (r *integrationRoutes) deleteIntegration(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_Response(c, http.StatusBadRequest, "ID must be integer")
		return
	}

	if err := r.uc.Delete(id); err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (r *integrationRoutes) handleRedirect(c *gin.Context) {
	code := c.Query("code")
	clientID := c.Query("client_id")
	if code == "" {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Authorization code is required"})
		return
	}

	tokens, err := r.uc.GetTokensByAuthCode(code, clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	integr, err := r.uc.GetIntegrationByClientID(clientID)
	if err != nil {
		var integration *entity.Integration
		integration.ClientID = clientID
		integration.AuthCode = code
		integration.Token = tokens

		if err := r.uc.Create(integration); err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
			return
		}
	}
	integr.Token = tokens

	if err := r.uc.Create(integr); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"integration": integr,
		"tokens":      tokens,
	})
}
