package amocrm

import (
	"strconv"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

//APIContactsResponse структура ответа от api
type APIContactsResponseDTO struct {
	Page  int `json:"_page"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Embedded struct {
		Contacts []struct {
			ID                 int         `json:"id"`
			Name               string      `json:"name"`
			FirstName          string      `json:"first_name"`
			LastName           string      `json:"last_name"`
			ResponsibleUserID  int         `json:"responsible_user_id"`
			GroupID            int         `json:"group_id"`
			CreatedBy          int         `json:"created_by"`
			UpdatedBy          int         `json:"updated_by"`
			CreatedAt          int         `json:"created_at"`
			UpdatedAt          int         `json:"updated_at"`
			ClosestTaskAt      interface{} `json:"closest_task_at"`
			IsDeleted          bool        `json:"is_deleted"`
			IsUnsorted         bool        `json:"is_unsorted"`
			CustomFieldsValues []struct {
				FieldID   int    `json:"field_id"`
				FieldName string `json:"field_name"`
				FieldCode string `json:"field_code"`
				FieldType string `json:"field_type"`
				Values    []struct {
					Value    string `json:"value"`
					EnumID   int    `json:"enum_id"`
					EnumCode string `json:"enum_code"`
				} `json:"values"`
			} `json:"custom_fields_values"`
			AccountID int `json:"account_id"`
			Links     struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"_links"`
			Embedded struct {
				Tags      []interface{} `json:"tags"`
				Companies []struct {
					ID    int `json:"id"`
					Links struct {
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
					} `json:"_links"`
				} `json:"companies"`
			} `json:"_embedded"`
		} `json:"contacts"`
	} `json:"_embedded"`
}

// ProcessContacts преобразует API-ответ в структуры Contact
func (r *APIContactsResponseDTO) ProcessContacts() *[]entity.Contact {
	contacts := make([]entity.Contact, 0, len(r.Embedded.Contacts))
	id := 0

	for _, apiContact := range r.Embedded.Contacts {
		contact := entity.Contact{
			ID:     id,
			Name:   apiContact.Name,
			Status: "unsync",
		}

		var email string
		for _, field := range apiContact.CustomFieldsValues {
			if field.FieldCode == "EMAIL" && len(field.Values) > 0 {
				email = field.Values[0].Value
				break
			}
		}

		if email == "" {
			contact.Status = "invalid"
		} else {
			contact.Email = email
		}

		contacts = append(contacts, contact)
		id++
	}

	return &contacts
}

type WebhookContactDTO struct {
	Contacts struct {
		Add []struct {
			ID                string `json:"id"`
			Name              string `json:"name"`
			ResponsibleUserID string `json:"responsible_user_id"`
			DateCreate        string `json:"date_create"`
			LastModified      string `json:"last_modified"`
			CreatedUserID     string `json:"created_user_id"`
			ModifiedUserID    string `json:"modified_user_id"`
			CompanyName       string `json:"company_name"`
			LinkedCompanyID   string `json:"linked_company_id"`
			AccountID         string `json:"account_id"`
			CustomFields      []struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Values []struct {
					Value string `json:"value"`
					Enum  string `json:"enum"`
				} `json:"values"`
				Code string `json:"code"`
			} `json:"custom_fields"`
			LinkedLeadsID map[string]struct {
				ID string `json:"ID"`
			} `json:"linked_leads_id"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
			Type      string `json:"type"`
		} `json:"add"`
	} `json:"contacts"`
}

// ConvertWebhookToGlobalContacts преобразует данные вебхука в срез GlobalContact
func ConvertWebhookToGlobalContacts(webhookData WebhookContactDTO) ([]entity.GlobalContact, error) {
	globalContacts := make([]entity.GlobalContact, 0, len(webhookData.Contacts.Add))

	for _, apiContact := range webhookData.Contacts.Add {
		accountID, _ := strconv.Atoi(apiContact.AccountID)

		globalContact := entity.GlobalContact{
			AccountID: accountID, // Преобразуем строку в int
			Name:      apiContact.Name,
			Status:    "unsync", // По умолчанию, статус "unsync"
		}

		var email string
		for _, field := range apiContact.CustomFields {
			if field.Code == "EMAIL" && len(field.Values) > 0 {
				email = field.Values[0].Value
			}
		}

		if email != "" {
			globalContact.Email = email
		} else {
			globalContact.Status = "invalid"
		}

		globalContacts = append(globalContacts, globalContact)
	}

	return globalContacts, nil
}
