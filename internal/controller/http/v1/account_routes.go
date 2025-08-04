package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/dto"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"

	"github.com/gin-gonic/gin"
)

type accountRoutes struct {
	uc account.AccountUseCase
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewAccountRoutes(handler *gin.RouterGroup, uc account.AccountUseCase) {
	r := &accountRoutes{uc}

	h := handler.Group("/accounts")
	{
		h.POST("/", r.createAccount)
		h.GET("/", r.getAccounts)
		h.GET("/:id", r.getAccount)
		h.GET("/:id/integrations", r.getAccountIntegrations)
		h.PUT("/:id", r.updateAccount)
		h.DELETE("/:id", r.deleteAccount)
		h.GET("/contacts", r.getContacts)
	}
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

func (r *oauthRoutes) GetContacts(token string) (*dto.ContactsResponse, error) {
	fullURL := r.makeRouteURL("/api/v4/contacts")

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	body, err := sendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while sending request to get contacts")
	}

	var apiResponse dto.APIContactsResponse

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return apiResponse.ToContactsResponse(), nil

}

//шлюз на внутренние методы
func (r *accountRoutes) createAccount(c *gin.Context) {
	var account entity.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		error_Response(c, http.StatusBadRequest, err.Error())
		return
	}

	account.CreatedAt = int(time.Now().Unix())

	if err := r.uc.Create(&account); err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, account)
}

//шлюз на методы
func (r *accountRoutes) getAccounts(c *gin.Context) {
	accounts, err := r.uc.GetAccounts()
	if err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}

//шлюз на методы
func (r *accountRoutes) getAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_Response(c, http.StatusBadRequest, "invalid account ID")
		return
	}

	account, err := r.uc.GetAccount(id)
	if err != nil {
		error_Response(c, http.StatusNotFound, "account not found")
		return
	}

	c.JSON(http.StatusOK, account)
}

//шлюз на методы
func (r *accountRoutes) getAccountIntegrations(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_Response(c, http.StatusBadRequest, "invalid account ID")
		return
	}

	integration, err := r.uc.GetAccountIntegrations(accountID)
	if err != nil {
		error_Response(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, integration)
}

func (r *accountRoutes) updateAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_Response(c, http.StatusBadRequest, "invalid account ID")
		return
	}

	var account entity.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		error_Response(c, http.StatusBadRequest, err.Error())
		return
	}

	account.ID = id
	if err := r.uc.Update(&account); err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, account)
}

func (r *accountRoutes) deleteAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_Response(c, http.StatusBadRequest, "invalid account ID")
		return
	}

	if err := r.uc.Delete(id); err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func error_Response(c *gin.Context, code int, err string) {
	c.AbortWithStatusJSON(code, errorResponse{Error: err})
}

const (
	REFRESH_THRESHOLD_SEC = 3600
)

//GetTokens - возвращает токены аккаунта по id
//UpdateTokens - обновляет токены пользователю (активному)
//MainUser - выбор активного пользователя. map[int]bool - int id аккаунта, bool - актив или нет
func (r *accountRoutes) StartTokenRefresher() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			tokens, err := r.uc.GetTokens()
			if err != nil || tokens == nil {
				continue
			}

			expiryTime := tokens.ServerTime + tokens.ExpiresIn
			if expiryTime-int(time.Now().Unix()) <= REFRESH_THRESHOLD_SEC {
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

func (r *accountRoutes) updateToken(c *gin.Context) {
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

func (r *accountRoutes) UpdateTokens(client_id string, config *config.Config) (*entity.Token, error) {
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
	base, err := url.Parse(config.BaseUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	base.Path = path.Join(base.Path, "/oauth2/access_token")
	fullURL := base.String()
	return r.sendTokenRequest(data, fullURL)
}
