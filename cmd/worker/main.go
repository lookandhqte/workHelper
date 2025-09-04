package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
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
	workAmount, err := strconv.Atoi(cfg.WorkerAmount)
	if err != nil {
		fmt.Println(err)
	}
	workers := []*worker.Worker{}
	for i := 0; i < workAmount; i++ {
		w := worker.NewWorker(cfg.BeanstalkAddr, *accountUC)
		workers = append(workers, w)
		go w.Start()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	for i := 0; i < workAmount; i++ {
		workers[i].Stop()
		fmt.Printf("worker %v stopped..\n", i)
	}
}
