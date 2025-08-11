package v1

import (
	"net/http"
	"strconv"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/provider"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/integration"
	"github.com/gin-gonic/gin"
)

//integrationRoutes роутер для интеграций
type integrationRoutes struct {
	uc       integration.UseCase
	provider provider.Provider
}

const (
	//BaseURL сайт
	BaseURL = "https://spetser.amocrm.ru/"
)

//NewIntegrationRoutes создает роуты для /integrations
func NewIntegrationRoutes(handler *gin.RouterGroup, uc integration.UseCase, provider provider.Provider) {
	r := &integrationRoutes{uc: uc, provider: provider}

	h := handler.Group("/integrations")
	{
		h.POST("/", r.createIntegration)
		h.GET("/", r.getIntegrations)
		h.PUT("/:id", r.updateIntegration)
		h.DELETE("/:id", r.deleteIntegration)
	}
}

//createIntegration создает интеграцию
func (r *integrationRoutes) createIntegration(c *gin.Context) {
	var integration entity.Integration

	if err := c.ShouldBindJSON(&integration); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.Create(&integration); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, integration)
}

//getIntegrations возвращает интеграции
func (r *integrationRoutes) getIntegrations(c *gin.Context) {
	integrations, err := r.uc.ReturnAll()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, integrations)
}

//updateIntegration обновляет интеграцию
func (r *integrationRoutes) updateIntegration(c *gin.Context) {
	var integration entity.Integration
	if err := c.ShouldBindJSON(&integration); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.Update(&integration); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, integration)
}

//deleteIntegration удаляет интеграцию
func (r *integrationRoutes) deleteIntegration(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "ID must be integer")
		return
	}

	if err := r.uc.Delete(id); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
