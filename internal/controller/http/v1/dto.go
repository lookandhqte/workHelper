package v1

import (
	entity "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

type APIUnisenderRequestDTO struct {
	UnisenderKey string `json:"unisender_key"`
	AccountID    int    `json:"account_id"`
}

// ConvertToGlobalContacts преобразует []Contact в []GlobalContact
func ConvertToGlobalContacts(contacts *[]entity.Contact) []entity.GlobalContact {
	globalContacts := make([]entity.GlobalContact, 0, len(*contacts))

	for _, contact := range *contacts {
		globalContact := entity.GlobalContact{
			AccountID: contact.AccountID,
			Email:     contact.Email,
			Name:      contact.Name,
			Status:    contact.Status,
		}
		globalContacts = append(globalContacts, globalContact)
	}

	return globalContacts
}

// type WebhookContactDTO struct {
// 	Contacts struct {
// 		Add []struct {
// 			ID                string `json:"id"`
// 			Name              string `json:"name"`
// 			ResponsibleUserID string `json:"responsible_user_id"`
// 			DateCreate        string `json:"date_create"`
// 			LastModified      string `json:"last_modified"`
// 			CreatedUserID     string `json:"created_user_id"`
// 			ModifiedUserID    string `json:"modified_user_id"`
// 			CompanyName       string `json:"company_name"`
// 			LinkedCompanyID   string `json:"linked_company_id"`
// 			AccountID         string `json:"account_id"`
// 			CustomFields      []struct {
// 				ID     string `json:"id"`
// 				Name   string `json:"name"`
// 				Values []struct {
// 					Value string `json:"value"`
// 					Enum  string `json:"enum"`
// 				} `json:"values"`
// 				Code string `json:"code"`
// 			} `json:"custom_fields"`
// 			LinkedLeadsID map[string]struct {
// 				ID string `json:"ID"`
// 			} `json:"linked_leads_id"`
// 			CreatedAt string `json:"created_at"`
// 			UpdatedAt string `json:"updated_at"`
// 			Type      string `json:"type"`
// 		} `json:"add"`
// 	} `json:"contacts"`
// }

// // ConvertWebhookToGlobalContacts преобразует данные вебхука в срез GlobalContact
// func ConvertWebhookToGlobalContacts(webhookData WebhookContactDTO) ([]entity.GlobalContact, error) {
// 	globalContacts := make([]entity.GlobalContact, 0, len(webhookData.Contacts.Add))

// 	for _, apiContact := range webhookData.Contacts.Add {
// 		accountID, _ := strconv.Atoi(apiContact.AccountID)

// 		globalContact := entity.GlobalContact{
// 			AccountID: accountID, // Преобразуем строку в int
// 			Name:      apiContact.Name,
// 			Status:    "unsync", // По умолчанию, статус "unsync"
// 		}

// 		var email string
// 		for _, field := range apiContact.CustomFields {
// 			if field.Code == "EMAIL" && len(field.Values) > 0 {
// 				email = field.Values[0].Value
// 			}
// 		}

// 		if email != "" {
// 			globalContact.Email = email
// 		} else {
// 			globalContact.Status = "invalid"
// 		}

// 		globalContacts = append(globalContacts, globalContact)
// 	}

// 	return globalContacts, nil
// }
