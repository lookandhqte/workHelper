package app

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	controllerhttp "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/controller/http"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/provider"
	accountUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"
	storageUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/storage"
	"github.com/gin-gonic/gin"
)

type dependencies struct {
	cfg       *config.Config
	AccountUC *accountUC.UseCase
	//TasksProducer *producer.TaskProducer
	Provider *provider.Provider
}

func composeDependencies() *dependencies {
	cfg := config.Load()

	storage := storageUC.NewStorage(cfg)
	return &dependencies{
		cfg:       cfg,
		AccountUC: accountUC.New(*storage),
		Provider:  provider.New(),
		//TasksProducer: producer.NewTaskProducer(cfg.BeanstalkAddr),
	}
}

func setupRouter(deps *dependencies) *gin.Engine {
	router := gin.Default()

	controllerhttp.NewRouter(router, *deps.AccountUC, *deps.Provider)

	return router
}
