package app

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	controllerhttp "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/controller/http"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/repo/persistent"
	accountUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"
	integrationUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/integration"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/pkg/cache"

	"github.com/gin-gonic/gin"
)

type dependencies struct {
	cfg           *config.Config
	AccountUC     *accountUC.AccountUseCase
	IntegrationUC *integrationUC.IntegrationUseCase
}

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

func setupRouter(deps *dependencies) *gin.Engine {
	router := gin.Default()

	controllerhttp.NewRouter(router, *deps.AccountUC, *deps.IntegrationUC)

	return router
}
