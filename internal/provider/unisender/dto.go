package unisender

type UnisenderResponseDTO struct {
	Result struct {
		ListID       int `json:"listId"`
		SearchParams struct {
			TagIds string `json:"tagIds"`
			Type   string `json:"type"`
		} `json:"searchParams"`
		Count string `json:"count"`
	} `json:"result"`
}

type ListUnisenderDTO struct {
	Result []struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"result"`
}
