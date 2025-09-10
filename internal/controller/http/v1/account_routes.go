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
		h.GET("/hh_ids", r.getVacanciesDescriptions)
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
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	account, err := r.uc.Return()
	if err != nil {
		log.Printf("error while getting account: %v", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	token.CreatedAt = int(time.Now().Unix())
	token.AccountID = account.ID
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

func (r *accountRoutes) getVacanciesDescriptions(c *gin.Context) {
	account, err := r.uc.Return()
	if err != nil {
		log.Printf("error while getting account: %v", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses, err := r.provider.HH.ReturnResponses(account.Token.AccessToken)
	if err != nil {
		log.Printf("error while getting responses: %v", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	resps := (*responses)[:1]

	prompts, err := r.provider.HH.ReturnPromptData(account.Token.AccessToken, &resps)
	if err != nil {
		log.Printf("error while getting prompt data: %v", err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	for _, prompt := range *prompts {
		message, err := r.provider.DeepSeek.GetVacancySoprovod(prompt)

		if err != nil {
			log.Printf("error while getting soprovod: %v", err)
			errorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		for i := range resps {
			if prompt.ID == resps[i].VacancyID {
				resps[i].Message = message
				err := r.provider.HH.DoResponse(account.Token.AccessToken, &resps[i])
				if err != nil {
					fmt.Printf("err while do response: %v\n", err)
				}
				break
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"responses": resps,
	})
}
