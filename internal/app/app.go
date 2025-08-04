package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	server          *http.Server
	backgroundTasks []BackgroundTask
}

func New() *App {
	return &App{}
}

type BackgroundTask func(ctx context.Context)

func (a *App) AddBackgroundTask(task BackgroundTask) {
	a.backgroundTasks = append(a.backgroundTasks, task)
}

// Run запускает приложение
func (a *App) Run() {

	deps := composeDependencies()
	router := setupRouter(deps)
	a.server = &http.Server{
		Addr:    deps.cfg.HTTPAddr,
		Handler: router,
	}

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
