package deepseek

// responseDTO ответ дипсика
type responseDTO struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}


