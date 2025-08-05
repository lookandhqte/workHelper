package http

import (
	"net/http"

	v1 "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/controller/http/v1"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/integration"

	"github.com/gin-gonic/gin"
)

type Router struct {
	accountUC     account.AccountUseCase
	integrationUC integration.IntegrationUseCase
	client        *http.Client
}

func NewRouter(
	r *gin.Engine,
	accountUC account.AccountUseCase,
	integrationUC integration.IntegrationUseCase,
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
func (r *Router) accountRoutes(api *gin.RouterGroup) {
	v1.NewAccountRoutes(api, r.accountUC)
}

func (r *Router) integrationRoutes(api *gin.RouterGroup) {
	v1.NewIntegrationRoutes(api, r.integrationUC, r.client)
}
