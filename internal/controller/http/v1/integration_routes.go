package v1

import (
	"amocrm_golang/internal/entity"
	"amocrm_golang/internal/usecase/integration"
	"net/http"
	"strconv"

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
	}
}

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

func (r *integrationRoutes) getIntegrations(c *gin.Context) {
	integrations, err := r.uc.Return(nil) // Используем Return вместо GetAll
	if err != nil {
		error_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, integrations)
}

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
