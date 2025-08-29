package app

import (
	"github.com/gin-gonic/gin"
	"github.com/lookandhqte/workHelper/config"
	controllerhttp "github.com/lookandhqte/workHelper/internal/controller/http"
	"github.com/lookandhqte/workHelper/internal/provider"
	accountUC "github.com/lookandhqte/workHelper/internal/usecase/account"
	storageUC "github.com/lookandhqte/workHelper/internal/usecase/storage"
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
