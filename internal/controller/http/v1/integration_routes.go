package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/integration"

	"github.com/gin-gonic/gin"
)

type integrationRoutes struct {
	uc integration.IntegrationUseCase
}

const (
	BASE_URL = "https://spetser.amocrm.ru/"
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

	contacts, err := r.uc.GetContacts(integr.Token.AccessToken)
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

	tokens, err := r.GetTokensByAuthCode(code, clientID)
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

func (m *integrationRoutes) GetTokensByAuthCode(code string, client_id string) (*entity.Token, error) {

	integration, err := m.uc.GetIntegrationByClientID(client_id)
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
