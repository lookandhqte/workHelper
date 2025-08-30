package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	entity "github.com/lookandhqte/workHelper/internal/entity"
	"github.com/lookandhqte/workHelper/internal/producer"
	"github.com/lookandhqte/workHelper/internal/provider"
	accountUC "github.com/lookandhqte/workHelper/internal/usecase/account"
)

// accountRoutes роутер для аккаунта
type accountRoutes struct {
	uc       accountUC.UseCase
	producer producer.TaskProducer
	provider provider.Provider
}

// NewAccountRoutes создает роуты для /account
func NewAccountRoutes(handler *gin.RouterGroup, producer producer.TaskProducer, provider provider.Provider, uc accountUC.UseCase) {
	r := &accountRoutes{producer: producer, provider: provider, uc: uc}

	h := handler.Group("/account")
	{
		h.POST("/", r.createAccount)
		h.GET("/", r.getAccount)
		h.PUT("/", r.updateAccount)
		h.DELETE("/", r.deleteAccount)
		h.GET("/redirect", r.handleRedirect)
	}
}

// createAccount создает акаунт
func (r *accountRoutes) createAccount(c *gin.Context) {
	var account entity.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	account.CreatedAt = int(time.Now().Unix())
	payload, err := json.Marshal(account)
	if err != nil {
		log.Printf("err while marshal: %v", err)
	}
	task := &entity.Task{
		Payload: payload,
		Type:    "account_creating",
	}

	if err := r.producer.CreateTask(task); err != nil {
		log.Printf("create task failed: %v\n", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, account)
}

// getAccount возвращает аккаунт
func (r *accountRoutes) getAccount(c *gin.Context) {

	account, err := r.uc.Return()
	if err != nil {
		errorResponse(c, http.StatusNotFound, "account not found")
		return
	}

	c.JSON(http.StatusOK, account)
}

// updateAccount обновляет аккаунт
func (r *accountRoutes) updateAccount(c *gin.Context) {

	var account entity.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	payload, err := json.Marshal(account)
	if err != nil {
		log.Printf("create task failed: %v\n", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	task := &entity.Task{
		Payload: payload,
		Type:    "account_updating",
	}

	if err := r.producer.CreateTask(task); err != nil {
		log.Printf("create task failed: %v\n", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, account)
}

// deleteAccount удаляет аккаунт
func (r *accountRoutes) deleteAccount(c *gin.Context) {

	if err := r.uc.Delete(); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// errorResponse ответ с ошибкой
func errorResponse(c *gin.Context, code int, err string) {
	c.AbortWithStatusJSON(code, fmt.Errorf("error: %v", err).Error())
}

// handleRedirect редирект
func (r *accountRoutes) handleRedirect(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		errorResponse(c, http.StatusBadRequest, "authorization code is required")
		return
	}

	token, err := r.provider.HH.GetToken(code)
	if err != nil {
		log.Printf("err while get token func handle redirect:  %v\n", err)
	}

	account, err := r.uc.Return()
	if err != nil {
		log.Printf("error while getting account: %v", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	account.Token = *token

	payload, err := json.Marshal(account)
	if err != nil {
		log.Printf("create task failed: %v\n", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	task := &entity.Task{
		Payload: payload,
		Type:    "account_updating",
	}

	if err := r.producer.CreateTask(task); err != nil {
		log.Printf("create task failed: %v\n", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"account": account,
	})
}
