package app

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	controllerhttp "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/controller/http"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/producer"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/provider"
	accountUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"
	contactsUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/contacts"
	integrationUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/integration"
	storageUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/storage"
	cache "git.amocrm.ru/gelzhuravleva/amocrm_golang/pkg/cache"
	"github.com/gin-gonic/gin"
)

type dependencies struct {
	cfg           *config.Config
	AccountUC     *accountUC.UseCase
	IntegrationUC *integrationUC.UseCase
	ContactsUC    *contactsUC.UseCase
	TasksProducer *producer.TaskProducer
	Provider      *provider.Provider
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
		ContactsUC:    contactsUC.New(*storage),
		TasksProducer: producer.NewTaskProducer(cfg.BeanstalkAddr),
		Provider:      provider.New(),
	}
}

func setupRouter(deps *dependencies) *gin.Engine {
	router := gin.Default()

	controllerhttp.NewRouter(router, *deps.AccountUC, *deps.IntegrationUC, *deps.TasksProducer, *deps.Provider, *deps.ContactsUC)

	return router
}
