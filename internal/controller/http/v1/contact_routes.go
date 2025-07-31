package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/token"

	"github.com/gin-gonic/gin"
)

type contactRoutes struct {
	uc token.TokenUseCase
}

func NewContactRoutes(handler *gin.RouterGroup, uc token.TokenUseCase) {
	r := &contactRoutes{uc}

	h := handler.Group("/contacts")
	{
		h.GET("/", r.getContacts)
	}
}

func (r *contactRoutes) getContacts(c *gin.Context) {
	tokens, err := r.uc.GetTokens()
	if err != nil {
		err_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	if tokens == nil {
		err_Response(c, http.StatusUnauthorized, "authentication required")
		return
	}

	if tokens.AccessToken == "" {
		err_Response(c, http.StatusUnauthorized, "access token missing")
		return
	}

	contacts, err := r.GetContacts(tokens.AccessToken)
	if err != nil {
		err_Response(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"contacts": contacts,
	})

}

func (r *contactRoutes) GetContacts(token string) (*entity.Contacts, error) {
	cfg := config.Load()
	base, err := url.Parse(cfg.BaseUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	base.Path = path.Join(base.Path, "/api/v4/contacts")
	fullURL := base.String()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var apiResponse struct {
		Embedded struct {
			Contacts []struct {
				Name               string `json:"name"`
				CustomFieldsValues []struct {
					FieldCode string `json:"field_code"`
					Values    []struct {
						Value string `json:"value"`
					} `json:"values"`
				} `json:"custom_fields_values"`
			} `json:"contacts"`
		} `json:"_embedded"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	var contacts entity.Contacts
	for _, contact := range apiResponse.Embedded.Contacts {
		sc := entity.Contact{
			Name: contact.Name,
		}

		for _, field := range contact.CustomFieldsValues {
			if field.FieldCode == "EMAIL" && len(field.Values) > 0 {
				email := field.Values[0].Value
				sc.Email = &email
				break
			}
		}

		contacts = append(contacts, sc)
	}

	return &contacts, nil

}
