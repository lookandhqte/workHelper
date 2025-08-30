package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lookandhqte/workHelper/config"
	"github.com/lookandhqte/workHelper/internal/usecase/account"
	"github.com/lookandhqte/workHelper/internal/usecase/storage"
	"github.com/lookandhqte/workHelper/internal/worker"
)

func main() {
	cfg := config.Load()
	storage := storage.NewStorage(cfg.StorageType, cfg.WorkerDSN)
	accountUC := account.New(*storage)
	w := worker.NewWorker(cfg.BeanstalkAddr, *accountUC)

	go w.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	w.Stop()
}
