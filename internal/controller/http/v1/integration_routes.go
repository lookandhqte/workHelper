package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/integration"

	"github.com/gin-gonic/gin"
)

type integrationRoutes struct {
	uc integration.IntegrationUseCase
}

func NewIntegrationRoutes(handler *gin.RouterGroup, uc integration.IntegrationUseCase) {
	r := &integrationRoutes{uc}

	h := handler.Group("/integrations")
	{
		h.POST("/", r.createIntegration)
		h.GET("/", r.getIntegrations)
		h.PUT("/:id", r.updateIntegration)
		h.DELETE("/:id", r.deleteIntegration)
		h.GET("/redirect", r.handleRedirect)
	}
}

//шлюз на внутренние методы
func (r *integrationRoutes) createIntegration(c *gin.Context) {
	var integration entity.Integration
	if err := c.ShouldBindJSON(&integration); err != nil {
		error_Response(c, http.StatusBadRequest, err.Error())
		return
	}

	if integration.AccountID == 0 {
		error_Response(c, http.StatusBadRequest, "account ID is required")
		return
	}

	if err := r.uc.Create(&integration); err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, integration)
}

//шлюз на внутренние методы
func (r *integrationRoutes) getIntegrations(c *gin.Context) {
	integrations, err := r.uc.Return(nil)
	if err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, integrations)
}

//шлюз на внутренние методы
func (r *integrationRoutes) updateIntegration(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_Response(c, http.StatusBadRequest, "ID must be integer")
		return
	}

	var integration entity.Integration
	if err := c.ShouldBindJSON(&integration); err != nil {
		error_Response(c, http.StatusBadRequest, err.Error())
		return
	}

	integration.AccountID = id
	if err := r.uc.Update(&integration); err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, integration)
}

//шлюз на внутренние методы
func (r *integrationRoutes) deleteIntegration(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_Response(c, http.StatusBadRequest, "ID must be integer")
		return
	}

	if err := r.uc.Delete(id); err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (r *integrationRoutes) getTokens(c *gin.Context) {
	tokens, err := r.uc.GetTokensByAuthCode(c.Query("code"), c.Query("client_id"))

	if err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Printf("Tokens:\n%s\n%s", tokens.AccessToken, tokens.RefreshToken)

	//В этом моменте желательно передать токены на ожидание через гофунку модели accounts, там обновлять токены сразу как придут новые в активном аккаунте.
	r.uc.CreateTokens(tokens)
}

func (r *integrationRoutes) handleRedirect(c *gin.Context) {
	code := c.Query("code")
	clientID := c.Query("client_id")
	if code == "" {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Authorization code is required"})
		return
	}

	tokens, err := r.uc.GetTokensByAuthCode(code, clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	r.uc.CreateTokens(tokens)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"tokens": tokens,
	})
}
