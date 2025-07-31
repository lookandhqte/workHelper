package v1

import (
	"amocrm_golang/config"
	"amocrm_golang/internal/entity"
	"amocrm_golang/internal/usecase/token"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

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
}

func (r *oauthRoutes) handleRedirect(c *gin.Context) {
	// Извлекаем параметры из запроса
	code := c.DefaultQuery("code", "")
	clientID := c.DefaultQuery("client_id", "")

	if clientID != "" {
		err := os.Setenv("CLIENT_ID", clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errResponse{Error: "Failed to set CLIENT_ID"})
			return
		}
	}

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
	data.Set("client_secret", config.ClientSecret) //config.ClientSecret
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", config.RedirectURI) //config.RedirectURI

	req, err := http.NewRequest("POST", config.AccessTokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error while making request: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while requesting: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response: %v", err)
	}

	responseData := &entity.Token{}

	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}

	return responseData, nil
}

func (r *oauthRoutes) getTokens(c *gin.Context) {
	cfg := config.Load()
	tokens, err := r.GetTokensByAuthCode(c.Param("code"), c.Param("client_id"), cfg)

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

func (r *oauthRoutes) updateToken(c *gin.Context) {
	cfg := config.Load()
	id := c.Param("client_id")
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
		return nil, fmt.Errorf("error while getting refresh token")
	}
	data := url.Values{}
	data.Set("client_id", client_id)
	data.Set("client_secret", config.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refresh)
	data.Set("redirect_uri", config.RedirectURI)

	req, err := http.NewRequest("POST", config.AccessTokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	responseData := &entity.Token{}

	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return responseData, nil
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
