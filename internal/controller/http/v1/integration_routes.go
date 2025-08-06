package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/integration"

	"github.com/gin-gonic/gin"
)

type integrationRoutes struct {
	uc     integration.IntegrationUseCase
	client *http.Client
}

const (
	BASE_URL          = "https://spetser.amocrm.ru/"
	REF_THRESHOLD_SEC = 3600
)

func NewIntegrationRoutes(handler *gin.RouterGroup, uc integration.IntegrationUseCase, client *http.Client) {
	r := &integrationRoutes{uc: uc, client: client}

	h := handler.Group("/integrations")
	{
		h.POST("/", r.createIntegration)
		h.GET("/", r.getIntegrations)
		h.PUT("/:account_id", r.updateIntegration)
		h.DELETE("/:account_id", r.deleteIntegration)
		h.GET("/redirect", r.handleRedirect)
		h.GET("/:account_id/contacts", r.getContacts)
		h.POST("/:account_id/refresh", r.needToRef)
	}
}

//Проблемы:
//контакты не сохраняются конкретной интеграции (по логике они после этого должны высасываться аккаунтом)
//токены не шифруются
func (r *integrationRoutes) needToRef(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("account_id"))
	integration, _ := r.uc.GetIntegration(id)
	newTokens, err := r.UpdateTokens(integration.ClientID)
	if err != nil {
		log.Printf("[Acc:%d] Failed to refresh token: %v", integration.AccountID, err)
		return
	}
	//Нужно будет добавить шифрование токенов из auth.
	//добавить на этом этапе шифрование токенов
	integration.Token = newTokens

	if err := r.uc.Update(integration); err != nil {
		log.Printf("[Acc:%d] Failed to save tokens: %v", integration.AccountID, err)
	}
}
func (r *integrationRoutes) getContacts(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("client_id"))
	id++

	integration, err := r.uc.GetIntegration(id)

	if err != nil {
		error_Response(c, http.StatusUnauthorized, "error in  get contacts func -> get int by client id")
		return
	}

	if integration.Token.AccessToken == "" {
		error_Response(c, http.StatusUnauthorized, "access token missing")
		return
	}

	contacts, err := r.GetContacts(integration.Token.AccessToken)
	if err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"contacts": contacts,
	})
}

func (r *integrationRoutes) createIntegration(c *gin.Context) {
	var integration entity.Integration
	if err := c.ShouldBindJSON(&integration); err != nil {
		error_Response(c, http.StatusBadRequest, err.Error())
		return
	}
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
	integrations, err := r.uc.Return()
	if err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, integrations)
}

func (r *integrationRoutes) updateIntegration(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("account_id"))
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

func (r *integrationRoutes) deleteIntegration(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("account_id"))
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

	tokens, err := r.GetTokensByAuthCode(code, clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	integration, err := r.uc.ReturnByClientID(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	integration.Token = tokens

	if err := r.uc.Update(integration); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"integration": integration,
		"tokens":      tokens,
	})
}

func (r *integrationRoutes) PrepareData(datacase string, integration *entity.Integration, code string, client_id string, contacts *ContactsResponse) url.Values {
	data := url.Values{}
	switch datacase {
	case "authorization_code":
		data.Set("client_id", client_id)
		data.Set("client_secret", integration.SecretKey)
		data.Set("grant_type", "authorization_code")
		data.Set("code", code)
		data.Set("redirect_uri", integration.RedirectUrl)
	case "refresh_token":
		data.Set("client_id", client_id)
		data.Set("client_secret", integration.SecretKey)
		data.Set("grant_type", "refresh_token")
		data.Set("refresh_token", integration.Token.RefreshToken)
		data.Set("redirect_uri", integration.RedirectUrl)
	}
	return data
}

func (r *integrationRoutes) GetTokensByAuthCode(code string, client_id string) (*entity.Token, error) {

	integration, err := r.uc.ReturnByClientID(client_id)
	if err != nil {
		return nil, fmt.Errorf("error n func get integr by client id -> error in get tokens method")
	}

	data := r.PrepareData("authorization_code", integration, code, client_id, nil)

	fullurl := r.MakeRouteURL("/oauth2/access_token")
	req, err := r.PreparePostRequest(fullurl, data)
	if err != nil {
		return nil, err
	}
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

	responseData := &entity.Token{}
	if err := json.Unmarshal(body, responseData); err != nil {
		return nil, err
	}

	return responseData, nil
}

func (r *integrationRoutes) UpdateTokens(client_id string) (*entity.Token, error) {
	integration, err := r.uc.ReturnByClientID(client_id)
	if err != nil {
		fmt.Print("error in func update tokens -> return by client id")
	}

	data := r.PrepareData("refresh_token", integration, "", client_id, nil)

	fullUrl := r.MakeRouteURL(BASE_URL)
	return r.SendTokenRequest(data, fullUrl)
}

func (r *integrationRoutes) GetContacts(token string) (*ContactsResponse, error) {
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

	var apiResponse APIContactsResponse

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return apiResponse.ToContactsResponse(), nil

}

func (r *integrationRoutes) MakeRouteURL(pathi string) string {
	base, _ := url.Parse(BASE_URL)

	base.Path = path.Join(base.Path, pathi)
	fullURL := base.String()
	return fullURL
}

func (r *integrationRoutes) PreparePostRequest(url string, data url.Values) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (r *integrationRoutes) SendRequest(req *http.Request) ([]byte, error) {
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

func (r *integrationRoutes) ParseTokenResponse(body []byte) (*entity.Token, error) {
	responseData := &entity.Token{}
	if err := json.Unmarshal(body, responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

func (r *integrationRoutes) SendTokenRequest(data url.Values, url string) (*entity.Token, error) {
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
