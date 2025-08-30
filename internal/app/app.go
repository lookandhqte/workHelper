package app

import (
	"context"
	"errors"
	"log"
	"net/http"
)

// App структура приложения
type App struct {
	server *http.Server
}

const (
	ShutdownTime = 5
)

// New возвращает экземпляр приложения App
func New() *App {
	return &App{}
}

// Run запускает экземпляр App
func (a *App) Run() {

	deps := composeDependencies()
	router := setupRouter(deps)

	a.server = &http.Server{
		Addr:    deps.cfg.HTTPAddr,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on %s", deps.cfg.HTTPAddr)
		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server error: %v", err)
		}
	}()

}

// Stop останавливает экземпляр App
func (a *App) Stop() {
	defer log.Println("Server exited. All tasks completed")

	log.Println("Shutting down server...")
	_, cancel := context.WithTimeout(context.Background(), ShutdownTime)
	defer cancel()

}
