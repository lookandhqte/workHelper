package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type App struct {
	server      *http.Server
	taskChan    chan func(ctx context.Context)
	wg          sync.WaitGroup
	shutdownCtx context.Context
	cancel      context.CancelFunc
}

func New() *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		taskChan:    make(chan func(ctx context.Context), 10),
		shutdownCtx: ctx,
		cancel:      cancel,
	}
}

func (a *App) AddTask(task func(ctx context.Context)) {
	select {
	case a.taskChan <- task:
		a.wg.Add(1)
	default:
		log.Println("Task channel is full, droppng task")
	}
}

func (a *App) StartTaskWorker() {
	go func() {
		for {
			select {
			case task := <-a.taskChan:
				go func(t func(ctx context.Context)) {
					defer a.wg.Done()
					t(a.shutdownCtx)
				}(task)
			case <-a.shutdownCtx.Done():
				return
			}
		}
	}()
}

func (a *App) Run() {

	a.StartTaskWorker()

	deps := composeDependencies()
	router := setupRouter(deps)

	a.server = &http.Server{
		Addr:    deps.cfg.HTTPAddr,
		Handler: router,
	}

	a.AddTask(func(ctx context.Context) {
		log.Printf("Server starting on %s", deps.cfg.HTTPAddr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	})

	a.AddTask(deps.IntegrationUC.Start(&a.wg))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	a.cancel()

	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All background tasks completed")
	case <-time.After(5 * time.Second):
		log.Println("Timeout waiting for tasks to complete")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Server exited gracefully")
}
