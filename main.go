package main

import (
	"amocrm_golang/database"
	"amocrm_golang/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация хранилища и сервисов
	storage := database.NewMemoryStorage()
	accountService := database.NewAccountService(storage)
	integrationService := database.NewIntegrationService(storage)

	// Создаем gin-роутер
	r := gin.Default()

	// Настройка роутов
	routes.SetupAccountRoutes(r, accountService, storage)
	routes.SetupIntegrationRoutes(r, integrationService)

	// Запускаем сервер на порту 2020
	r.Run(":2020")
}
