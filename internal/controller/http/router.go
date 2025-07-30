package http

import (
	v1 "amocrm_golang/internal/controller/http/v1"
	"amocrm_golang/internal/usecase/account"
	"amocrm_golang/internal/usecase/integration"
	"amocrm_golang/internal/usecase/token"

	"github.com/gin-gonic/gin"
)

type Router struct {
	accountUC     account.AccountUseCase
	integrationUC integration.IntegrationUseCase
	tokenUC       token.TokenUseCase
}

func NewRouter(
	r *gin.Engine,
	accountUC account.AccountUseCase,
	integrationUC integration.IntegrationUseCase,
	tokenUC token.TokenUseCase,
) {
	router := &Router{
		accountUC:     accountUC,
		integrationUC: integrationUC,
		tokenUC:       tokenUC,
	}

	api := r.Group("/v1")
	{
		router.accountRoutes(api)
		router.integrationRoutes(api)
		router.oauthRoutes(api)
	}
}

func (r *Router) accountRoutes(api *gin.RouterGroup) {
	v1.NewAccountRoutes(api, r.accountUC)
}

func (r *Router) integrationRoutes(api *gin.RouterGroup) {
	v1.NewIntegrationRoutes(api, r.integrationUC)
}

func (r *Router) oauthRoutes(api *gin.RouterGroup) {
	v1.NewTokenRoutes(api, r.tokenUC)
}
