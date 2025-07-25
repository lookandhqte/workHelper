package app

import (
	"amocrm_golang/config"
	controllerhttp "amocrm_golang/internal/controller/http"
	"amocrm_golang/internal/repo/persistent"
	"amocrm_golang/internal/usecase"
	"amocrm_golang/pkg/cache"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	fils "github.com/swaggo/files"
	sw "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func Run() {
	// Инициализация конфига
	cfg := config.Load()

	// Инициализация зависимостей
	memoryCache := cache.NewCache()
	storage := persistent.NewMemoryStorage(memoryCache)

	accountUseCase := usecase.NewAccountUseCase(storage)
	integrationUseCase := usecase.NewIntegrationUseCase(storage)

	// Настройка HTTP сервера
	router := gin.Default()
	router.GET("/swagger/*any", sw.WrapHandler(fils.Handler))

	controllerhttp.NewRouter(router, *accountUseCase, *integrationUseCase)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
