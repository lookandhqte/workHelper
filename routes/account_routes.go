package routes

import (
	"amocrm_golang/auth"
	"amocrm_golang/database"
	"amocrm_golang/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

			// ID будет сгенерирован в accountService.CreateAccount()
			accessToken, err := auth.GenerateJWT(account.ID, auth.AccessTokenExpiry)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
				return
			}

			refreshToken, err := auth.GenerateJWT(account.ID, auth.RefreshTokenExpiry)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
				return
			}

			account.AccessToken = accessToken
			account.RefreshToken = refreshToken
			account.CreatedAt = time.Now()
			account.TokenExpires = time.Now().Add(auth.AccessTokenExpiry)
			account.Expires = 7 // для кэширования

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
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
				return
			}

			account, err := storage.GetAccountWithCache(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Проверяем срок действия токена
			if time.Now().After(account.TokenExpires) {
				newToken, err := auth.GenerateJWT(account.ID, auth.AccessTokenExpiry)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to refresh token"})
					return
				}
				account.AccessToken = newToken
				account.TokenExpires = time.Now().Add(auth.AccessTokenExpiry)
				storage.UpdateAccount(account)
			}

			// Проверяем expires (дни кэширования)
			if account.Expires <= 0 {
				account.Expires = 7
				storage.UpdateAccount(account)
			}

			c.JSON(http.StatusOK, account)
		})

		accountGroup.GET("/:id/integrations", func(c *gin.Context) {
			accountID, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID format"})
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
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
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
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
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
