package v1

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	entity "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

type APIUnisenderRequestDTO struct {
	UnisenderKey string `json:"unisender_key"`
	AccountID    int    `json:"account_id"`
}

// ConvertToGlobalContacts преобразует []Contact в []GlobalContact
func ConvertToGlobalContactsDTO(contacts *[]entity.Contact) []entity.GlobalContact {
	globalContacts := make([]entity.GlobalContact, 0, len(*contacts))

	for _, contact := range *contacts {
		globalContact := entity.GlobalContact{
			AccountID: contact.AccountID,
			Email:     contact.Email,
			Status:    contact.Status,
		}
		globalContacts = append(globalContacts, globalContact)
	}

	return globalContacts
}

func ConvertWebhookToGlobalContactsDTO(formData map[string][]string) *entity.GlobalContact {
	globalContact := &entity.GlobalContact{}
	if ids, ok := formData["account[id]"]; ok && len(ids) > 0 {
		id, err := strconv.Atoi(ids[0])
		if err != nil {
			log.Printf("error while strconv account id func convert to global from webhook: %v", err)
		}
		globalContact.AccountID = id
	}

	if ids, ok := formData["contacts[add][0][id]"]; ok && len(ids) > 0 {
		id, err := strconv.Atoi(ids[0])
		if err != nil {
			log.Printf("error while strconv contact id func convert to global from webhook: %v", err)
		}
		globalContact.ID = id
	}
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	for key, values := range formData {
		if strings.Contains(key, "[code]") && len(values) > 0 && values[0] == "EMAIL" {
			prefix := strings.Split(key, "[code]")[0]
			if emailValues, ok := formData[prefix+"[values][0][value]"]; ok && len(emailValues) > 0 {
				if emailRegex.MatchString(emailValues[0]) {
					globalContact.Email = emailValues[0]
				}
			}
		}
	}

	globalContact.Status = "unsync"

	return globalContact
}
