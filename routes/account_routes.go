package routes

import (
	"amocrm_golang/database"
	"amocrm_golang/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetupAccountRoutes(r *gin.Engine, accountService *database.AccountService, storage *database.MemoryStorage) {
	accountGroup := r.Group("/accounts")
	{
		accountGroup.POST("/", func(c *gin.Context) {
			var account model.Account
			if err := c.ShouldBindJSON(&account); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			account.ID = uuid.New()
			if err := accountService.CreateAccount(&account); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, account)
		})

		accountGroup.GET("/", func(c *gin.Context) {
			accounts, err := accountService.GetAccountList()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, accounts)
		})

		accountGroup.GET("/:id", func(c *gin.Context) {
			id, err := uuid.Parse(c.Param("id"))
			account, err := accountService.GetAccountByID(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, account)
		})

		accountGroup.GET("/:id/integrations", func(c *gin.Context) {
			accountID, err := uuid.Parse(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
				return
			}

			integration, err := storage.GetAccountIntegrations(accountID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, integration)
		})

		accountGroup.PUT("/:id", func(c *gin.Context) {
			id, err := uuid.Parse(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
				return
			}

			var account model.Account
			if err := c.ShouldBindJSON(&account); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			account.ID = id
			if err := accountService.UpdateAccount(&account); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, account)
		})

		accountGroup.DELETE("/:id", func(c *gin.Context) {
			id, err := uuid.Parse(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
				return
			}

			if err := accountService.DeleteAccount(id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "account deleted"})
		})
	}
}
