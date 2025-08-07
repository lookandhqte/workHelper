package app

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	controllerhttp "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/controller/http"
	accountUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"
	integrationUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/integration"
	storageUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/storage"
	cache "git.amocrm.ru/gelzhuravleva/amocrm_golang/pkg/cache"
	"github.com/gin-gonic/gin"
)

type dependencies struct {
	cfg           *config.Config
	AccountUC     *accountUC.UseCase
	IntegrationUC *integrationUC.UseCase
}

func composeDependencies() *dependencies {
	cfg := config.Load()

	memoryCache := cache.NewCache()

	if cfg.StorageType == "database" {
		memoryCache = nil
	}

	storage := storageUC.NewStorage(memoryCache, cfg)

	return &dependencies{
		cfg:           cfg,
		AccountUC:     accountUC.New(*storage),
		IntegrationUC: integrationUC.New(*storage),
	}
}

func setupRouter(deps *dependencies) *gin.Engine {
	router := gin.Default()

	controllerhttp.NewRouter(router, *deps.AccountUC, *deps.IntegrationUC)

	return router
}
