package http

import (
	v1 "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/controller/http/v1"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/producer"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/provider"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/contacts"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/integration"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/worker"
	"github.com/gin-gonic/gin"
)

//Router абстракция
type Router struct {
	accountUC     account.UseCase
	integrationUC integration.UseCase
	contactsUC    contacts.UseCase
	provider      provider.Provider
	producer      producer.TaskProducer
	workers       worker.TaskWorkers
}

//NewRouter создает новый роутер
func NewRouter(
	r *gin.Engine,
	accountUC account.UseCase,
	integrationUC integration.UseCase,
	contactsUC contacts.UseCase,
	producer producer.TaskProducer,
	provider provider.Provider,
	workers worker.TaskWorkers,
) {
	router := &Router{
		accountUC:     accountUC,
		integrationUC: integrationUC,
		provider:      provider,
		producer:      producer,
		contactsUC:    contactsUC,
		workers:       workers,
	}

	api := r.Group("/v1")
	{
		router.accountRoutes(api)
		router.integrationRoutes(api)
		router.contactsRoutes(api)
	}
}

//accountRoutes создает роуты для аккаунта
func (r *Router) accountRoutes(api *gin.RouterGroup) {
	v1.NewAccountRoutes(api, r.accountUC, r.provider)
}

//accountRoutes создает роуты для аккаунта
func (r *Router) contactsRoutes(api *gin.RouterGroup) {
	v1.NewContactsRoutes(api, r.contactsUC, r.producer, &r.workers)
}

//integrationRoutes создает роуты для интеграций
func (r *Router) integrationRoutes(api *gin.RouterGroup) {
	v1.NewIntegrationRoutes(api, r.integrationUC, r.provider)
}
