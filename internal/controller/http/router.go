package http

import (
	v1 "amocrm_golang/internal/controller/http/v1"
	"amocrm_golang/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Router struct {
	accountUC     usecase.AccountUseCase
	integrationUC usecase.IntegrationUseCase
}

func NewRouter(
	r *gin.Engine,
	accountUC usecase.AccountUseCase,
	integrationUC usecase.IntegrationUseCase,
) {
	router := &Router{
		accountUC:     accountUC,
		integrationUC: integrationUC,
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
	v1.NewIntegrationRoutes(api, r.integrationUC)
}
