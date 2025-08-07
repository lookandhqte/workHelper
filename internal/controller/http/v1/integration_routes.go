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
		h.GET("/:id/contacts", r.getContacts)
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
	tokens, err := r.uc.GetTokens(integration.TokenID)
	if err != nil {
		log.Printf("err while getting tokens: %v", err)
	}
	tokens.UnisenderKey = unisenderKey

	if err := r.uc.UpdateToken(tokens); err != nil {
		log.Printf("err while updating tokens: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"tokens": tokens,
	})
}

//needToRef обновляет токены при необходимости
func (r *integrationRoutes) needToRef(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	integration, _ := r.uc.ReturnOne(id)
	newTokens, err := r.UpdateTokens(integration.ClientID)

	if err != nil {
		log.Printf("[Acc:%d] Failed to refresh token: %v", integration.AccountID, err)
		return
	}

	if err := r.uc.UpdateToken(newTokens); err != nil {
		log.Printf("Failed to save tokens: %v", err)
	}
	integration.TokenID = newTokens.AccountID

	if err := r.uc.Update(integration); err != nil {
		log.Printf("[Acc:%d] Failed to save tokens: %v", integration.AccountID, err)
	}
}

//getContacts возвращает контакты
func (r *integrationRoutes) getContacts(c *gin.Context) {
	fmt.Printf("gin context: %v", c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("error in func get contacts: %v", err)
	}

	integration, err := r.uc.ReturnOne(id)

	if err != nil {
		fmt.Printf("error in func return one from intfunc get contacts: %v", err)
	}

	tokens, err := r.uc.GetTokens(integration.TokenID)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "error in  get tokens func -> get int by client id")
		return
	}
	if tokens.AccessToken == "" {
		errorResponse(c, http.StatusUnauthorized, "access token missing")
		return
	}

	contacts, err := r.GetContacts(tokens.AccessToken)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"contacts": contacts,
	})
}

//createIntegration создает интеграцию
func (r *integrationRoutes) createIntegration(c *gin.Context) {
	var integration entity.Integration

	if err := c.ShouldBindJSON(&integration); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if integration.AccountID == 0 {
		errorResponse(c, http.StatusBadRequest, "account ID is required")
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "ID must be integer")
		return
	}

	var integration entity.Integration
	if err := c.ShouldBindJSON(&integration); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	integration.AccountID = id
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
	fmt.Println(tokens)

	if err := r.uc.UpdateToken(tokens); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	integration, err := r.uc.ReturnByClientID(clientID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	integration.TokenID = tokens.ID
	integration.AuthCode = code

	if err := r.uc.Update(integration); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println(integration)

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
		tokens, err := r.uc.GetTokens(integration.TokenID)
		if err != nil {
			log.Printf("error case switch tokens: %v", err)
		}
		data.Set("refresh_token", tokens.RefreshToken)
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

	fullURL := r.MakeRouteURL("/oauth2/access_token")
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

	return responseData, nil
}

//UpdateTokens обновляет токены
func (r *integrationRoutes) UpdateTokens(clientID string) (*entity.Token, error) {
	integration, err := r.uc.ReturnByClientID(clientID)
	if err != nil {
		fmt.Print("error in func update tokens -> return by client id")
	}

	data := r.PrepareData("refresh_token", integration, "", clientID)

	fullURL := r.MakeRouteURL(BaseURL)
	return r.SendTokenRequest(data, fullURL)
}

//GetContacts возвращает контакты
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

//MakeRouteURL возвращает полный URL адрес
func (r *integrationRoutes) MakeRouteURL(pathi string) string {
	base, _ := url.Parse(BaseURL)

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
