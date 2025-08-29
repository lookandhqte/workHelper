package http

import (
	v1 "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/controller/http/v1"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/provider"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"
	"github.com/gin-gonic/gin"
)

// Router абстракция
type Router struct {
	accountUC account.UseCase
	provider  provider.Provider
	//producer  producer.TaskProducer
}

// NewRouter создает новый роутер
func NewRouter(
	r *gin.Engine,
	accountUC account.UseCase,
	//producer producer.TaskProducer,
	provider provider.Provider,
) {
	router := &Router{
		accountUC: accountUC,
		provider:  provider,
		//producer:  producer,
	}

	api := r.Group("/v1")
	{
		router.accountRoutes(api)
	}
}

// accountRoutes создает роуты для аккаунта
func (r *Router) accountRoutes(api *gin.RouterGroup) {
	v1.NewAccountRoutes(api, r.accountUC, r.provider)
}
