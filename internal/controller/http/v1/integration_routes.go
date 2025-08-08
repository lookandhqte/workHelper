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

//integrationRoutes роутер для интеграций
type integrationRoutes struct {
	uc     integration.UseCase
	client *http.Client
}

const (
	//BaseURL сайт
	BaseURL = "https://spetser.amocrm.ru/"
)

//NewIntegrationRoutes создает роуты для /integrations
func NewIntegrationRoutes(handler *gin.RouterGroup, uc integration.UseCase, client *http.Client) {
	r := &integrationRoutes{uc: uc, client: client}

	h := handler.Group("/integrations")
	{
		h.POST("/", r.createIntegration)
		h.GET("/", r.getIntegrations)
		h.PUT("/:id", r.updateIntegration)
		h.DELETE("/:id", r.deleteIntegration)
		h.GET("/redirect", r.handleRedirect)
		h.POST("/:id/refresh", r.needToRef)
		h.POST("/:id/unisender", r.saveUnisenderToken)
	}
}

//saveUnisenderToken сохраняет интеграции токен unisender
func (r *integrationRoutes) saveUnisenderToken(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("err while converting br: %v", err)
	}
	var request APIUnisenderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	unisenderKey := request.UnisenderKey
	integration, err := r.uc.ReturnOne(id)
	if err != nil {
		log.Printf("error while getting integration: %v", err)
	}
	integration.Token.UnisenderKey = unisenderKey

	if err := r.uc.Update(integration); err != nil {
		log.Printf("err while updating tokens: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"unisender key": unisenderKey,
	})
}

//needToRef обновляет токены при необходимости
func (r *integrationRoutes) needToRef(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("failed to parse id: %v", err)
		return
	}
	integration, _ := r.uc.ReturnOne(id)
	newTokens, err := r.UpdateTokens(integration.ClientID)

	if err != nil {
		log.Printf("[Acc:%d] Failed to refresh token: %v", integration.ID, err)
		return
	}

	integration.Token = newTokens

	if err := r.uc.Update(integration); err != nil {
		log.Printf("[Acc:%d] Failed to save tokens: %v", integration.ID, err)
	}
}

//createIntegration создает интеграцию
func (r *integrationRoutes) createIntegration(c *gin.Context) {
	var integration entity.Integration

	if err := c.ShouldBindJSON(&integration); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.Create(&integration); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, integration)
}

//getIntegrations возвращает интеграции
func (r *integrationRoutes) getIntegrations(c *gin.Context) {
	integrations, err := r.uc.ReturnAll()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, integrations)
}

//updateIntegration обновляет интеграцию
func (r *integrationRoutes) updateIntegration(c *gin.Context) {
	var integration entity.Integration
	if err := c.ShouldBindJSON(&integration); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.Update(&integration); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, integration)
}

//deleteIntegration удаляет интеграцию
func (r *integrationRoutes) deleteIntegration(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "ID must be integer")
		return
	}

	if err := r.uc.Delete(id); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

//handleRedirect редирект
func (r *integrationRoutes) handleRedirect(c *gin.Context) {
	code := c.Query("code")
	clientID := c.Query("client_id")
	if code == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("authorization code is required"))
		return
	}

	tokens, err := r.GetTokensByAuthCode(code, clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	integration, err := r.uc.ReturnByClientID(clientID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	integration.Token = tokens

	if err := r.uc.Update(integration); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"integration": integration,
		"tokens":      tokens,
	})
}

//PrepareData готовит url.Values для запроса
func (r *integrationRoutes) PrepareData(datacase string, integration *entity.Integration, code string, clientID string) url.Values {
	data := url.Values{}
	switch datacase {
	case "authorization_code":
		data.Set("client_id", clientID)
		data.Set("client_secret", integration.SecretKey)
		data.Set("grant_type", "authorization_code")
		data.Set("code", code)
		data.Set("redirect_uri", integration.RedirectURL)
	case "refresh_token":
		data.Set("client_id", clientID)
		data.Set("client_secret", integration.SecretKey)
		data.Set("grant_type", "refresh_token")
		data.Set("refresh_token", integration.Token.RefreshToken)
		data.Set("redirect_uri", integration.RedirectURL)
	}
	return data
}

//GetTokensByAuuthCode получает токены с кодом auth
func (r *integrationRoutes) GetTokensByAuthCode(code string, clientID string) (*entity.Token, error) {

	integration, err := r.uc.ReturnByClientID(clientID)
	if err != nil {
		return nil, fmt.Errorf("error n func get integr by client id -> error in get tokens method: %v", err)
	}

	data := r.PrepareData("authorization_code", integration, code, clientID)

	fullURL := MakeRouteURL("/oauth2/access_token", "")
	req, err := r.PreparePostRequest(fullURL, data)
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

	responseData.IntegrationID = integration.ID

	return responseData, nil
}

//UpdateTokens обновляет токены
func (r *integrationRoutes) UpdateTokens(clientID string) (*entity.Token, error) {
	integration, err := r.uc.ReturnByClientID(clientID)
	if err != nil {
		fmt.Print("error in func update tokens -> return by client id")
	}

	data := r.PrepareData("refresh_token", integration, "", clientID)

	fullURL := MakeRouteURL(BaseURL, "")
	return r.SendTokenRequest(data, fullURL)
}

//MakeRouteURL возвращает полный URL адрес
func MakeRouteURL(pathi string, baseURL string) string {
	base, _ := url.Parse(BaseURL)
	if baseURL != "" {
		base, _ = url.Parse(baseURL)
	}
	base.Path = path.Join(base.Path, pathi)
	fullURL := base.String()
	return fullURL
}

//PreparePostRequest готовит post запрос
func (r *integrationRoutes) PreparePostRequest(url string, data url.Values) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

//SendRequest отправляет запрос
func SendRequest(req *http.Request, r http.Client) ([]byte, error) {
	resp, err := r.Do(req)
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

//ParseTokenResponse парсит ответ в токены
func (r *integrationRoutes) ParseTokenResponse(body []byte) (*entity.Token, error) {
	responseData := &entity.Token{}
	if err := json.Unmarshal(body, responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

//SendTokenRequest отправляет запрос на получение токенов
func (r *integrationRoutes) SendTokenRequest(data url.Values, url string) (*entity.Token, error) {
	req, err := r.PreparePostRequest(url, data)
	if err != nil {
		return nil, fmt.Errorf("request preparation failed: %v", err)
	}

	body, err := SendRequest(req, *r.client)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	token, err := r.ParseTokenResponse(body)
	if err != nil {
		return nil, fmt.Errorf("response parsing failed: %v", err)
	}

	return token, nil
}
