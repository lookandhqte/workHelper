package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	entity "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	accountUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"
	"github.com/gin-gonic/gin"
)

//accountRoutes роутер для аккаунта
type accountRoutes struct {
	uc     accountUC.UseCase
	client *http.Client
}

const (
	SlicesCapacity = 10
)

//NewAccountRoutes создает роуты для /accounts
func NewAccountRoutes(handler *gin.RouterGroup, uc accountUC.UseCase, client *http.Client) {
	r := &accountRoutes{uc: uc, client: client}

	h := handler.Group("/accounts")
	{
		h.POST("/", r.createAccount)
		h.GET("/", r.getAccounts)
		h.GET("/:id", r.getAccount)
		h.GET("/:id/integrations", r.getAccountIntegrations)
		h.PUT("/:id", r.updateAccount)
		h.DELETE("/:id", r.deleteAccount)
		h.GET("/:id/contacts", r.getAccountContacts)
	}
}

//getAccountContacts возвращает контакты аккаунта по первой authed интеграции
func (r *accountRoutes) getAccountContacts(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("error while atoi id func get account contacts: %v", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	integrationsPtr, err := r.uc.ReturnIntegrations(id)
	if err != nil {
		log.Printf("error while getting account integrations: %v", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	integrations := *integrationsPtr
	idOfAuthedIntegration := 0
	for _, integration := range integrations {
		if integration.Token.AccessToken != "" {
			idOfAuthedIntegration = integration.ID
			break
		}
	}
	tokens := integrations[idOfAuthedIntegration].Token.AccessToken
	contactsResponse, err := r.GetContacts(tokens)
	if err != nil {
		log.Printf("error while getting contacts func get account contacts: %v", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	contacts := contactsResponse.ResponseToContacts(contactsResponse)

	account, err := r.uc.ReturnOne(integrations[0].AccountID)
	if err != nil {
		log.Printf("error while returning one account func get account contacts: %v", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	account.AccountContacts = *contacts

	err = r.uc.Update(account)

	if err != nil {
		log.Printf("error while updating account: %v", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"contacts":        contacts,
		"updated account": account,
	})

}

//GetContacts возвращает контакты
func (r *accountRoutes) GetContacts(token string) (*ContactsResponse, error) {
	fullURL := MakeRouteURL("/api/v4/contacts")

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	body, err := SendRequest(req, *r.client)
	if err != nil {
		return nil, fmt.Errorf("error while sending request to get contacts")
	}

	var apiResponse APIContactsResponse

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return apiResponse.ToContactsResponse(), nil

}

//createAccount создает акаунт
func (r *accountRoutes) createAccount(c *gin.Context) {
	var account entity.Account
	contacts := make([]entity.Contact, 0, SlicesCapacity)
	account.AccountContacts = contacts
	if err := c.ShouldBindJSON(&account); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	account.CreatedAt = int(time.Now().Unix())

	if err := r.uc.Create(&account); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, account)
}

//getAccounts возвращает аккаунты
func (r *accountRoutes) getAccounts(c *gin.Context) {
	accounts, err := r.uc.ReturnAll()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}

//getAccount возвращает аккаунт
func (r *accountRoutes) getAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid account ID")
		return
	}

	account, err := r.uc.ReturnOne(id)
	if err != nil {
		errorResponse(c, http.StatusNotFound, "account not found")
		return
	}

	c.JSON(http.StatusOK, account)
}

//getAccountIntegrations возвращает все интеграции аккаунта
func (r *accountRoutes) getAccountIntegrations(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid account ID")
		return
	}

	integration, err := r.uc.ReturnIntegrations(accountID)
	if err != nil {
		errorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, integration)
}

//updateAccount обновляет аккаунт
func (r *accountRoutes) updateAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid account ID")
		return
	}

	var account entity.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	account.ID = id
	if err := r.uc.Update(&account); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, account)
}

//deleteAccount удаляет аккаунт
func (r *accountRoutes) deleteAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid account ID")
		return
	}

	if err := r.uc.Delete(id); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

//errorResponse ответ с ошибкой
func errorResponse(c *gin.Context, code int, err string) {
	c.AbortWithStatusJSON(code, fmt.Errorf("error: %v", err).Error())
}
