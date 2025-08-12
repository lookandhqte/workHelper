package v1

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/producer"
	contactsUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/contacts"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/worker"
	"github.com/gin-gonic/gin"
)

//contactRoutes роутер для аккаунта
type contactsRoutes struct {
	uc           contactsUC.UseCase
	taskProducer producer.TaskProducer
	worker       worker.TaskWorker
}

//NewContactRoutes создает роуты для /contacts
func NewContactsRoutes(handler *gin.RouterGroup, uc contactsUC.UseCase, taskProducer producer.TaskProducer, worker worker.TaskWorker) {
	r := &contactsRoutes{uc: uc, taskProducer: taskProducer, worker: worker}

	h := handler.Group("/contacts")
	{
		h.POST("/", r.updateContacts)
	}
}

//updateAccount обновляет аккаунт
func (r *contactsRoutes) updateContacts(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		errorResponse(c, http.StatusInternalServerError, "Error reading body")
		return
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	bodyString := string(body)
	fmt.Println("Request body as string:", bodyString)

	data, err := url.ParseQuery(bodyString)
	if err != nil {
		log.Printf("Error parsing query: %v", err)
		errorResponse(c, http.StatusBadRequest, "Failed to parse query")
		return
	}

	contact := ConvertWebhookToGlobalContactsDTO(data)
	r.taskProducer.EnqueueCreateContactTask(contact)
	r.worker.ResolveCreateContactTask()
	if err := r.uc.Create(contact); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, contact)
}
