package app

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	controllerhttp "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/controller/http"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/repo/persistent"
	accountUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"
	integrationUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/integration"
	tokenUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/token"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/pkg/cache"

	"github.com/gin-gonic/gin"
)

// dependencies содержит все зависимости приложения
type dependencies struct {
	cfg           *config.Config
	AccountUC     *accountUC.AccountUseCase
	IntegrationUC *integrationUC.IntegrationUseCase
	TokenUC       *tokenUC.TokenUseCase
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
		TokenUC:       tokenUC.New(storage),
	}
}

// setupRouter настраивает маршруты приложения
func setupRouter(deps *dependencies) *gin.Engine {
	router := gin.Default()

	controllerhttp.NewRouter(router, *deps.AccountUC, *deps.IntegrationUC, *deps.TokenUC)

	return router
}
