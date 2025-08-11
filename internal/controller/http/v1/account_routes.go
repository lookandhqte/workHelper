package v1

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	entity "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	producer "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/producer"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/provider"
	accountUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"
	"github.com/gin-gonic/gin"
)

//accountRoutes роутер для аккаунта
type accountRoutes struct {
	uc           accountUC.UseCase
	provider     provider.Provider
	taskProducer producer.TaskProducer
}

//NewAccountRoutes создает роуты для /accounts
func NewAccountRoutes(handler *gin.RouterGroup, uc accountUC.UseCase, provider provider.Provider, taskProducer producer.TaskProducer) {
	r := &accountRoutes{uc: uc, provider: provider, taskProducer: taskProducer}

	h := handler.Group("/accounts")
	{
		h.POST("/", r.createAccount)
		h.GET("/", r.getAccounts)
		h.GET("/:id", r.getAccount)
		h.GET("/:id/integrations", r.getAccountIntegrations)
		h.PUT("/:id", r.updateAccount)
		h.DELETE("/:id", r.deleteAccount)
		h.GET(":id/redirect", r.handleRedirect)
		h.POST("/:id/refresh/:integration_id", r.needToRef)
	}
}

//createAccount создает акаунт
func (r *accountRoutes) createAccount(c *gin.Context) {
	var account entity.Account
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

func (r *accountRoutes) findIntegrationBy(findcase string, accountID int, integrationID int, clientID string) (*entity.Integration, error) {
	integrationsPtr, err := r.uc.ReturnIntegrations(accountID)
	integrations := *integrationsPtr
	if err != nil {
		log.Printf("err while return integrations func find integration by: %v", err)
	}
	if findcase == "" {
		return nil, fmt.Errorf("empty datacase")
	}
	switch findcase {
	case "byIntegrationID":
		for _, integration := range integrations {
			if integration.ID == integrationID {
				return &integration, nil
			}
		}
	case "byClientID":
		for _, integration := range integrations {
			if integration.ClientID == clientID {
				return &integration, nil
			}
		}
	case "nilnessToken":
		for _, integration := range integrations {
			if integration.Token != nil {
				return &integration, nil
			}
		}
	default:
		return nil, fmt.Errorf("no such integration")
	}
	return nil, fmt.Errorf("oh no my friend, no suchh datacase or integration.. : %v", err)
}

//needToRef обновляет токены при необходимости
func (r *accountRoutes) needToRef(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("failed to parse id: %v", err)
		return
	}
	integrationID, err := strconv.Atoi(c.Param("integration_id"))
	if err != nil {
		log.Printf("failed to parse int id: %v", err)
		return
	}
	account, _ := r.uc.ReturnOne(id)

	integration, err := r.findIntegrationBy("byIntegrationID", id, integrationID, "")
	if err != nil {
		log.Printf("failed find integration by -> need to ref: %v", err)
	}
	err = r.provider.Amo.UpdateTokens(integration)
	if err != nil {
		log.Printf("failed to update tokens func need to ref: %v", err)
	}

	if err := r.uc.Update(account); err != nil {
		log.Printf("[Acc:%d] Failed to save tokens: %v", err)
	}
}

//handleRedirect редирект
func (r *accountRoutes) handleRedirect(c *gin.Context) {

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("account id is required"))
		return
	}
	code := c.Query("code")
	clientID := c.Query("client_id")
	if code == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("authorization code is required"))
		return
	}
	integration, err := r.findIntegrationBy("byClientID", id, 0, clientID)
	if err != nil {
		log.Printf("failed to find integration func handle redirect: %v", err)
	}
	err = r.provider.Amo.GetTokens(integration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	account, err := r.uc.ReturnOne(id)
	if err != nil {
		log.Printf("error while getting account: %v", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	contacts, err := r.provider.Amo.GetContacts(integration.Token.AccessToken)
	if err != nil {
		log.Printf("error while getting account contacts in provider: %v", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	account.AccountContacts = *contacts

	if err := r.uc.Update(account); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	globalContacts := ConvertToGlobalContacts(contacts)
	if err := r.taskProducer.EnqueueSyncContactsTask(id, integration.ID, globalContacts); err != nil {
		log.Printf("enqueue task failed: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"account": account,
	})
}
