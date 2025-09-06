package http

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/lookandhqte/workHelper/internal/controller/http/v1"
	"github.com/lookandhqte/workHelper/internal/producer"
	"github.com/lookandhqte/workHelper/internal/provider"
	"github.com/lookandhqte/workHelper/internal/usecase/account"
)

// Router абстракция
type Router struct {
	provider  provider.Provider
	producer  producer.TaskProducer
	accountUC account.UseCase
}

// NewRouter создает новый роутер
func NewRouter(
	r *gin.Engine,
	producer producer.TaskProducer,
	provider provider.Provider,
	accountUC account.UseCase,
) {
	router := &Router{
		provider:  provider,
		producer:  producer,
		accountUC: accountUC,
	}

	api := r.Group("/v1")
	{
		router.accountRoutes(api)
	}
}

// accountRoutes создает роуты для аккаунта
func (r *Router) accountRoutes(api *gin.RouterGroup) {
	v1.NewAccountRoutes(api, r.producer, r.provider, r.accountUC)
}
