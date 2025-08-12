package producer

type WebhookContactDTO struct {
	Contacts struct {
		Add []struct {
			ID                string `form:"id"`
			Name              string `form:"name"`
			ResponsibleUserID string `form:"responsible_user_id"`
			DateCreate        string `form:"date_create"`
			LastModified      string `form:"last_modified"`
			CreatedUserID     string `form:"created_user_id"`
			ModifiedUserID    string `form:"modified_user_id"`
			CompanyName       string `form:"company_name"`
			LinkedCompanyID   string `form:"linked_company_id"`
			AccountID         string `form:"account_id"`
			CustomFields      []struct {
				ID     string `form:"id"`
				Name   string `form:"name"`
				Values []struct {
					Value string `form:"value"`
					Enum  string `form:"enum"`
				} `form:"values"`
				Code string `form:"code"`
			} `form:"custom_fields"`
			LinkedLeadsID struct {
				Num7551167 struct {
					ID string `form:"ID"`
				} `form:"7551167"`
			} `form:"linked_leads_id"`
			CreatedAt string `form:"created_at"`
			UpdatedAt string `form:"updated_at"`
			Type      string `form:"type"`
		} `form:"add"`
	} `form:"contacts"`
}
