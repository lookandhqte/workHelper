package deepseek

// responseDTO ответ дипсика
type responseDTO struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Response структура отклика
type ResponseToVacDTO struct {
	ResumeID  string `json:"resume_id"`
	VacancyID string `json:"vacancy_id"`
	Message   string `json:"message"` //максимум 10.000 символов
}
