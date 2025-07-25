package routes

import (
	"amocrm_golang/database"
	"amocrm_golang/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupIntegrationRoutes(r *gin.Engine, integrationService *database.IntegrationService) {
	integrationGroup := r.Group("/integrations")
	{
		integrationGroup.POST("/", func(c *gin.Context) {
			var integration model.Integration
			if err := c.ShouldBindJSON(&integration); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// AccountID должен приходить в теле запроса или генерироваться другим способом
			if integration.AccountID == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "account ID is required"})
				return
			}

			if err := integrationService.CreateIntegration(&integration); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, integration)
		})

		integrationGroup.GET("/", func(c *gin.Context) {
			integrations, err := integrationService.GetIntegrationList()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, integrations)
		})

		integrationGroup.PUT("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be integer"})
				return
			}

			var integration model.Integration
			if err := c.ShouldBindJSON(&integration); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			integration.AccountID = id
			if err := integrationService.UpdateIntegration(&integration); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, integration)
		})

		integrationGroup.DELETE("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be integer"})
				return
			}

			if err := integrationService.DeleteIntegration(id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "integration deleted"})
		})
	}
}
