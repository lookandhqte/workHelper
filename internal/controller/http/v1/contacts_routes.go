package v1

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/producer"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/provider"
	contactsUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/contacts"
	"github.com/gin-gonic/gin"
)

//contactRoutes роутер для аккаунта
type contactsRoutes struct {
	uc           contactsUC.UseCase
	provider     provider.Provider
	taskProducer producer.TaskProducer
}

//NewContactRoutes создает роуты для /contacts
func NewContactsRoutes(handler *gin.RouterGroup, uc contactsUC.UseCase, provider provider.Provider, taskProducer producer.TaskProducer) {
	r := &contactsRoutes{uc: uc, provider: provider, taskProducer: taskProducer}

	h := handler.Group("/contacts")
	{
		h.POST("/", r.updateContacts)
	}
}

func parseWebhookForm(form url.Values) producer.WebhookContactDTO {
	dto := producer.WebhookContactDTO{}

	for key, values := range form {
		if !strings.HasPrefix(key, "contacts[add][") {
			continue
		}

		parts := strings.Split(key, "[")
		if len(parts) < 4 {
			continue
		}

		contactIndex := strings.TrimSuffix(parts[2], "]")
		fieldPath := strings.TrimSuffix(strings.Join(parts[3:], "."), "]")

		if contactIndex != "0" {
			continue
		}

		value := values[0]
		switch fieldPath {
		case "id":
			dto.Contacts.Add[0].ID = value
		case "name":
			dto.Contacts.Add[0].Name = value
		case "account_id":
			dto.Contacts.Add[0].AccountID = value
		case "custom_fields":
			custom := parseCustomFields(form, contactIndex)
			dto.Contacts.Add[0].CustomFields = []struct {
				ID     string
				Code   string
				Values []struct{ Value string }
			}(custom)
		}
	}
	return dto
}

func parseCustomFields(form url.Values, contactIndex string) []struct {
	ID     string `form:"id"`
	Code   string `form:"code"`
	Values []struct {
		Value string `form:"value"`
	} `form:"values"`
} {
	fields := []struct {
		ID     string
		Code   string
		Values []struct{ Value string }
	}{}

	// Собираем все кастомные поля
	fieldIndex := 0
	for {
		baseKey := fmt.Sprintf("contacts[add][%s][custom_fields][%d]", contactIndex, fieldIndex)
		idKey := baseKey + "[id]"

		// Проверяем наличие поля
		if _, exists := form[idKey]; !exists {
			break
		}

		field := struct {
			ID     string
			Code   string
			Values []struct{ Value string }
		}{
			ID:   form.Get(idKey),
			Code: form.Get(baseKey + "[code]"),
		}

		// Собираем значения поля
		valueIndex := 0
		for {
			valueKey := fmt.Sprintf("%s[values][%d][value]", baseKey, valueIndex)
			if value := form.Get(valueKey); value != "" {
				field.Values = append(field.Values, struct{ Value string }{Value: value})
				valueIndex++
			} else {
				break
			}
		}

		fields = append(fields, field)
		fieldIndex++
	}

	return []struct {
		ID     string "form:\"id\""
		Code   string "form:\"code\""
		Values []struct {
			Value string "form:\"value\""
		} "form:\"values\""
	}(fields)
}

//updateAccount обновляет аккаунт
func (r *contactsRoutes) updateContacts(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		log.Printf("ParseForm error: %v", err)
		errorResponse(c, http.StatusBadRequest, "Invalid form data")
		return
	}

	request := parseWebhookForm(c.Request.Form)

	// Логируем результат парсинга для отладки
	log.Printf("Parsed webhook data: %+v", request)

	// Дальнейшая обработка...
	globalContacts, err := producer.ConvertWebhookToGlobalContacts(request)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.Update(globalContacts); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, globalContacts)
}
