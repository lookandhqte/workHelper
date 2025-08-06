package http

import (
	"net/http"

	v1 "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/controller/http/v1"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/integration"
	"github.com/gin-gonic/gin"
)

//Router абстракция
type Router struct {
	accountUC     account.UseCase
	integrationUC integration.UseCase
	client        *http.Client
}

//NewRouter создает новый роутер
func NewRouter(
	r *gin.Engine,
	accountUC account.UseCase,
	integrationUC integration.UseCase,
) {
	router := &Router{
		accountUC:     accountUC,
		integrationUC: integrationUC,
		client:        &http.Client{},
	}

	api := r.Group("/v1")
	{
		router.accountRoutes(api)
		router.integrationRoutes(api)
	}
}

//accountRoutes создает роуты для аккаунта
func (r *Router) accountRoutes(api *gin.RouterGroup) {
	v1.NewAccountRoutes(api, r.accountUC)
}

//integrationRoutes создает роуты для интеграций
func (r *Router) integrationRoutes(api *gin.RouterGroup) {
	v1.NewIntegrationRoutes(api, r.integrationUC, r.client)
}
