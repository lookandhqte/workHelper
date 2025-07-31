package http

import (
	v1 "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/controller/http/v1"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/contact"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/integration"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/token"

	"github.com/gin-gonic/gin"
)

type Router struct {
	accountUC     account.AccountUseCase
	integrationUC integration.IntegrationUseCase
	tokenUC       token.TokenUseCase
	contactUC     contact.ContactUseCase
}

func NewRouter(
	r *gin.Engine,
	accountUC account.AccountUseCase,
	integrationUC integration.IntegrationUseCase,
	tokenUC token.TokenUseCase,
	contactUC contact.ContactUseCase,
) {
	router := &Router{
		accountUC:     accountUC,
		integrationUC: integrationUC,
		tokenUC:       tokenUC,
		contactUC:     contactUC,
	}

	api := r.Group("/v1")
	{
		router.accountRoutes(api)
		router.integrationRoutes(api)
		router.oauthRoutes(api)
		router.contactRoutes(api)
	}
}

func (r *Router) contactRoutes(api *gin.RouterGroup) {
	v1.NewContactRoutes(api, r.contactUC)
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
