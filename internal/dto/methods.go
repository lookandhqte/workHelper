package dto

func (r *APIContactsResponse) ToContactsResponse() *ContactsResponse {
	contacts := make(ContactsResponse, 0, len(r.Embedded.Contacts))

	for _, contact := range r.Embedded.Contacts {
		cr := ContactResponse{
			Name: contact.Name,
		}

		for _, field := range contact.CustomFieldsValues {
			if field.FieldCode == "EMAIL" && len(field.Values) > 0 {
				email := field.Values[0].Value
				cr.Email = &email
				break
			}
		}

		contacts = append(contacts, cr)
	}

	return &contacts
}
