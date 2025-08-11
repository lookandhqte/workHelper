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
