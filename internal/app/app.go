package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strconv"
	"sync"
	"syscall"
	"time"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

//App структура приложения
type App struct {
	server      *http.Server
	taskChan    chan func(ctx context.Context)
	wg          sync.WaitGroup
	shutdownCtx context.Context
	cancel      context.CancelFunc
}

const (
	RefreshThreshold = 3600
	ShutdownTime     = 5
	SemaphorSize     = 10
)

func New() *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		taskChan:    make(chan func(ctx context.Context), SemaphorSize),
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

	a.StartGRPCServer(deps.AccountUC)

	a.server = &http.Server{
		Addr:    deps.cfg.HTTPAddr,
		Handler: router,
	}

	a.AddTask(func(ctx context.Context) {
		log.Printf("Server starting on %s", deps.cfg.HTTPAddr)
		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server error: %v", err)
		}
	})

	a.AddTask(func(ctx context.Context) {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				integrationsPtr, err := deps.IntegrationUC.ReturnAll()
				if err != nil {
					log.Printf("Failed to get active integrations: %v", err)
					return
				}
				integrations := *integrationsPtr
				sem := make(chan struct{}, SemaphorSize)
				for i := range integrations {
					a.wg.Add(1)
					sem <- struct{}{}

					go func(integration *entity.Integration) {
						defer a.wg.Done()
						defer func() { <-sem }()

						expiryTime := integration.Token.ServerTime + integration.Token.ExpiresIn
						now := time.Now().Unix()

						if expiryTime-int(now) <= RefreshThreshold {

							integrationID := strconv.Itoa(integration.ID)
							accountID := strconv.Itoa(integration.AccountID)
							base, _ := url.Parse("http://localhost:2020/v1/accounts/")
							base.Path = path.Join(base.Path, accountID)
							base.Path = path.Join(base.Path, "/refresh/")
							base.Path = path.Join(base.Path, integrationID)
							fullURL := base.String()
							req, _ := http.NewRequest(http.MethodPost, fullURL, nil)
							client := &http.Client{}
							resp, _ := client.Do(req)
							defer resp.Body.Close()
							if resp.StatusCode != http.StatusAccepted {
								log.Println("could not do req for updating refresh token")
							}
						}
					}(&integrations[i])
				}

				a.wg.Wait()
			case <-ctx.Done():
				log.Println("Token refresher stopped")
				return
			}
		}
	})

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
	case <-time.After(ShutdownTime * time.Second):
		log.Println("Timeout waiting for tasks to complete")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutdownTime*time.Second)
	defer cancel()
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Server exited gracefully")
}
