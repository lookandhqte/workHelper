package v1

//ContactResponse структура ответа
type ContactResponse struct {
	Name  string  `json:"name"`
	Email *string `json:"email"`
}

//ContactsResponse слайс сущностей ContactResponse
type ContactsResponse []ContactResponse

//APIContactsResponse структура ответа от api
type APIContactsResponse struct {
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

//ToContactsResponse превращает ответ api в ContactsResponse
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

type APIUnisenderRequest struct {
	UnisenderKey string `json:"unisender_key"`
	AccountID    int    `json:"account_id"`
}
