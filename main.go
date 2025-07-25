package main

import (
	"amocrm_golang/database"
	"amocrm_golang/routes"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	storage := database.NewMemoryStorage()
	accountService := database.NewAccountService(storage)
	integrationService := database.NewIntegrationService(storage)

	r := gin.Default()

	routes.SetupAccountRoutes(r, accountService, storage)
	routes.SetupIntegrationRoutes(r, integrationService)

	srv := &http.Server{
		Addr:    ":2020",
		Handler: r,
	}

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
