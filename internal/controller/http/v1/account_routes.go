package v1

import (
	"amocrm_golang/internal/entity"
	"amocrm_golang/internal/usecase/account"
	"amocrm_golang/pkg/auth"
	"net/http"
	"strconv"
	"time"

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

// @Summary     Create account
// @Description Create new account
// @ID          create-account
// @Tags  	    accounts
// @Accept      json
// @Produce     json
// @Param       request body entity.Account true "Account info"
// @Success     201 {object} entity.Account
// @Failure     400 {object} error_Response
// @Failure     500 {object} error_Response
// @Router      /v1/accounts [post]
func (r *accountRoutes) createAccount(c *gin.Context) {
	var account entity.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		error_Response(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, err := auth.GenerateJWT(account.ID, auth.AccessTokenExpiry)
	if err != nil {
		error_Response(c, http.StatusInternalServerError, "failed to generate access token")
		return
	}

	refreshToken, err := auth.GenerateJWT(account.ID, auth.RefreshTokenExpiry)
	if err != nil {
		error_Response(c, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}

	account.AccessToken = accessToken
	account.RefreshToken = refreshToken
	account.CreatedAt = time.Now()
	account.TokenExpires = time.Now().Add(auth.AccessTokenExpiry)
	account.Expires = 7

	if err := r.uc.Create(&account); err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, account)
}

// @Summary     List accounts
// @Description Get all accounts
// @ID          list-accounts
// @Tags  	    accounts
// @Accept      json
// @Produce     json
// @Success     200 {array} entity.Account
// @Failure     500 {object} error_Response
// @Router      /v1/accounts [get]
func (r *accountRoutes) getAccounts(c *gin.Context) {
	accounts, err := r.uc.GetAccounts()
	if err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// @Summary     Get account
// @Description Get account by ID
// @ID          get-account
// @Tags  	    accounts
// @Accept      json
// @Produce     json
// @Param       id path int true "Account ID"
// @Success     200 {object} entity.Account
// @Failure     400 {object} error_Response
// @Failure     404 {object} error_Response
// @Failure     500 {object} error_Response
// @Router      /v1/accounts/{id} [get]
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

// @Summary     Get account integrations
// @Description Get integrations for account
// @ID          get-account-integrations
// @Tags  	    accounts
// @Accept      json
// @Produce     json
// @Param       id path int true "Account ID"
// @Success     200 {object} entity.Integration
// @Failure     400 {object} error_Response
// @Failure     404 {object} error_Response
// @Failure     500 {object} error_Response
// @Router      /v1/accounts/{id}/integrations [get]
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

// @Summary     Update account
// @Description Update account
// @ID          update-account
// @Tags  	    accounts
// @Accept      json
// @Produce     json
// @Param       id path int true "Account ID"
// @Param       request body entity.Account true "Account info"
// @Success     200 {object} entity.Account
// @Failure     400 {object} error_Response
// @Failure     404 {object} error_Response
// @Failure     500 {object} error_Response
// @Router      /v1/accounts/{id} [put]
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

// @Summary     Delete account
// @Description Delete account
// @ID          delete-account
// @Tags  	    accounts
// @Accept      json
// @Produce     json
// @Param       id path int true "Account ID"
// @Success     204
// @Failure     400 {object} error_Response
// @Failure     404 {object} error_Response
// @Failure     500 {object} error_Response
// @Router      /v1/accounts/{id} [delete]
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
