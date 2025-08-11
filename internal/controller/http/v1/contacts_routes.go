package v1

import (
	"log"
	"net/http"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/producer"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/provider"
	contactsUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/contacts"
	"github.com/gin-gonic/gin"
)

//contactRoutes роутер для аккаунта
type contactsRoutes struct {
	uc           contactsUC.UseCase
	provider     provider.Provider
	taskProducer producer.TaskProducer
}

//NewContactRoutes создает роуты для /contacts
func NewContactsRoutes(handler *gin.RouterGroup, uc contactsUC.UseCase, provider provider.Provider, taskProducer producer.TaskProducer) {
	r := &contactsRoutes{uc: uc, provider: provider, taskProducer: taskProducer}

	h := handler.Group("/contacts")
	{
		h.POST("/", r.updateContacts)
	}
}

//updateAccount обновляет аккаунт
func (r *contactsRoutes) updateContacts(c *gin.Context) {

	request := &producer.WebhookContactDTO{}
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.taskProducer.EnqueueSyncWebhookContactsTask(*request); err != nil {
		log.Printf("enqueue task failed: %v", err)
	}

	globalContacts, err := producer.ConvertWebhookToGlobalContacts(*request)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := r.uc.Update(globalContacts); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, globalContacts)
}
