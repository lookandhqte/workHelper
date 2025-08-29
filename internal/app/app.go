package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// App структура приложения
type App struct {
	server    *http.Server
	wg        *sync.WaitGroup
	tasksChan chan func()
	stopChan  chan struct{}
	workers   int
}

const (
	ShutdownTime        = 5
	defaultWorkerAmount = 5
)

func New() *App {
	return &App{
		tasksChan: make(chan func(), 100),
		stopChan:  make(chan struct{}),
		wg:        &sync.WaitGroup{},
	}
}

func (a *App) StartWorkers() {
	for i := 0; i < a.workers; i++ {
		a.wg.Add(1)
		go a.worker()
	}
}

// func (a *App) StopWorkers() {
// 	for i:= 0; i< a.workers; i++ {
// 		a.worker()
// 	}
// }

func (a *App) worker() {
	defer a.wg.Done()

	for {
		select {
		case task := <-a.tasksChan:
			task()
		case <-a.stopChan:
			return
		}
	}
}

func (a *App) AddTask(task func()) {
	select {
	case a.tasksChan <- task:
	default:
		log.Println("Task channel full, task dropped")
	}
}

func (a *App) Run() {

	deps := composeDependencies()
	router := setupRouter(deps)

	a.server = &http.Server{
		Addr:    deps.cfg.HTTPAddr,
		Handler: router,
	}

	a.workers, _ = strconv.Atoi(deps.cfg.WorkerAmount)
	if a.workers == 0 {
		a.workers = defaultWorkerAmount
	}

	a.StartWorkers()

	go func() {
		log.Printf("Server starting on %s", deps.cfg.HTTPAddr)
		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server error: %v", err)
		}
	}()

}

func (a *App) Stop() {
	defer log.Println("Server exited. All tasks completed")

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTime)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Closing stop chan...")
	close(a.stopChan)

	log.Println("Waiting all tasks to complete...")
	//разобраться с wg wait
	a.wg.Wait()

	log.Println("Closing tasks chan...")
	close(a.tasksChan)
}
