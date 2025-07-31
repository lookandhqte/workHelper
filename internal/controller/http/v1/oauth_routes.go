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
	"time"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/token"

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
		h.GET("/contacts", r.getContacts)
	}

	r.StartTokenRefresher()
}

func (r *oauthRoutes) StartTokenRefresher() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			tokens, err := r.uc.GetTokens()
			if err != nil || tokens == nil {
				continue
			}

			expiryTime := tokens.ServerTime + tokens.ExpiresIn
			refresh_thr, err := r.uc.GetConst("refresh_threshold")
			if err != nil {
				fmt.Printf("you don't have consts refresh thr")
			}
			if expiryTime-int(time.Now().Unix()) <= refresh_thr {
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
	data.Set("client_secret", config.ClientSecret) //config.CllientSecret
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refresh)
	data.Set("redirect_uri", config.RedirectURI) //config.RedirectURI
	base, err := url.Parse(config.BaseUrl)       //config.BaseUrl
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	base.Path = path.Join(base.Path, "/oauth2/access_token")
	fullURL := base.String()
	return r.sendTokenRequest(data, fullURL)
}

func (r *oauthRoutes) handleRedirect(c *gin.Context) {
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
	base, err := url.Parse(config.BaseUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	base.Path = path.Join(base.Path, "/oauth2/access_token")
	fullURL := base.String()

	return r.sendTokenRequest(data, fullURL)
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
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(data.Encode()))
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

	if resp.StatusCode == http.StatusBadRequest {
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

func (r *oauthRoutes) getContacts(c *gin.Context) {
	tokens, err := r.uc.GetTokens()
	fmt.Println(tokens)
	if err != nil {
		err_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	if tokens == nil {
		err_Response(c, http.StatusUnauthorized, "authentication required")
		return
	}

	if tokens.AccessToken == "" {
		err_Response(c, http.StatusUnauthorized, "access token missing")
		return
	}

	contacts, err := r.GetContacts(tokens.AccessToken)
	if err != nil {
		err_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"contacts": contacts,
	})

}

func (r *oauthRoutes) GetContacts(token string) (*entity.Contacts, error) {
	cfg := config.Load()
	base, err := url.Parse(cfg.BaseUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	base.Path = path.Join(base.Path, "/api/v4/contacts")
	fullURL := base.String()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var apiResponse struct {
		Embedded struct {
			Contacts []struct {
				Name               string `json:"name"`
				CustomFieldsValues []struct {
					FieldCode string `json:"field_code"`
					Values    []struct {
						Value string `json:"value"`
					} `json:"values"`
				} `json:"custom_fields_values"`
			} `json:"contacts"`
		} `json:"_embedded"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	var contacts entity.Contacts
	for _, contact := range apiResponse.Embedded.Contacts {
		sc := entity.Contact{
			Name: contact.Name,
		}

		for _, field := range contact.CustomFieldsValues {
			if field.FieldCode == "EMAIL" && len(field.Values) > 0 {
				email := field.Values[0].Value
				sc.Email = &email
				break
			}
		}

		contacts = append(contacts, sc)
	}

	return &contacts, nil

}
