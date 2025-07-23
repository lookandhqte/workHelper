package routes

import (
	"amocrm_golang/database"
	"amocrm_golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

			integration.AccountID = uuid.New()
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
			id, err := uuid.Parse(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
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
			id, err := uuid.Parse(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
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
