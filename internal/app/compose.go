package app

import (
	"amocrm_golang/config"
	controllerhttp "amocrm_golang/internal/controller/http"
	"amocrm_golang/internal/repo/persistent"
	accountUC "amocrm_golang/internal/usecase/account"
	integrationUC "amocrm_golang/internal/usecase/integration"
	"amocrm_golang/pkg/cache"

	"github.com/gin-gonic/gin"
	fils "github.com/swaggo/files"
	sw "github.com/swaggo/gin-swagger"
)

// dependencies содержит все зависимости приложения
type dependencies struct {
	cfg           *config.Config
	AccountUC     *accountUC.AccountUseCase
	IntegrationUC *integrationUC.IntegrationUseCase
}

// composeDependencies инициализирует все зависимости
func composeDependencies() *dependencies {
	cfg := config.Load()
	memoryCache := cache.NewCache()
	storage := persistent.NewMemoryStorage(memoryCache)

	return &dependencies{
		cfg:           cfg,
		AccountUC:     accountUC.New(storage),
		IntegrationUC: integrationUC.New(storage),
	}
}

// setupRouter настраивает маршруты приложения
func setupRouter(deps *dependencies) *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", sw.WrapHandler(fils.Handler))

	controllerhttp.NewRouter(router, *deps.AccountUC, *deps.IntegrationUC)

	return router
}
