package app

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	controllerhttp "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/controller/http"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/producer"
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
}

func composeDependencies() *dependencies {
	cfg := config.Load()

	memoryCache := cache.NewCache()

	if cfg.StorageType == "database" {
		memoryCache = nil
	}

	storage := storageUC.NewStorage(memoryCache, cfg)
	tasksProducer := producer.NewTaskProducer(cfg.BeanstalkAddr)
	return &dependencies{
		cfg:           cfg,
		AccountUC:     accountUC.New(*storage),
		IntegrationUC: integrationUC.New(*storage),
		ContactsUC:    contactsUC.New(*storage),
		TasksProducer: tasksProducer,
	}
}

func setupRouter(deps *dependencies) *gin.Engine {
	router := gin.Default()

	controllerhttp.NewRouter(router, *deps.AccountUC, *deps.IntegrationUC, *deps.TasksProducer)

	return router
}
