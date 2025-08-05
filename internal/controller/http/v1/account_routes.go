package v1

import (
	"net/http"
	"strconv"
	"time"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"

	"github.com/gin-gonic/gin"
)

type accountRoutes struct {
	uc account.AccountUseCase
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewAccountRoutes(handler *gin.RouterGroup, uc account.AccountUseCase) {
	r := &accountRoutes{uc}

	h := handler.Group("/accounts")
	{
		h.POST("/", r.createAccount)
		h.GET("/", r.getAccounts)
		h.GET("/:id", r.getAccount)
		h.GET("/:id/integrations", r.getAccountIntegrations)
		h.PUT("/:id", r.updateAccount)
		h.DELETE("/:id", r.deleteAccount)
	}
}

func (r *accountRoutes) createAccount(c *gin.Context) {
	var account entity.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		error_Response(c, http.StatusBadRequest, err.Error())
		return
	}

	account.CreatedAt = int(time.Now().Unix())

	if err := r.uc.Create(&account); err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, account)
}
func (r *accountRoutes) getAccounts(c *gin.Context) {
	accounts, err := r.uc.GetAccounts()
	if err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (r *accountRoutes) getAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_Response(c, http.StatusBadRequest, "invalid account ID")
		return
	}

	account, err := r.uc.GetAccount(id)
	if err != nil {
		error_Response(c, http.StatusNotFound, "account not found")
		return
	}

	c.JSON(http.StatusOK, account)
}

func (r *accountRoutes) getAccountIntegrations(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_Response(c, http.StatusBadRequest, "invalid account ID")
		return
	}

	integration, err := r.uc.GetAccountIntegrations(accountID)
	if err != nil {
		error_Response(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, integration)
}

func (r *accountRoutes) updateAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_Response(c, http.StatusBadRequest, "invalid account ID")
		return
	}

	var account entity.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		error_Response(c, http.StatusBadRequest, err.Error())
		return
	}

	account.ID = id
	if err := r.uc.Update(&account); err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, account)
}

func (r *accountRoutes) deleteAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_Response(c, http.StatusBadRequest, "invalid account ID")
		return
	}

	if err := r.uc.Delete(id); err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func error_Response(c *gin.Context, code int, err string) {
	c.AbortWithStatusJSON(code, errorResponse{Error: err})
}
