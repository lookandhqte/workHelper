package v1

import (
	"fmt"

	"github.com/lookandhqte/workHelper/internal/provider/hh"
)

// Response структура отклика
type Response struct {
	VacancyURL string `json:"url"`
	ResumeID   string `json:"resume_id"`
	VacancyID  string `json:"vacancy_id"`
	Message    string `json:"message"` //максимум 10.000 символов
}

func resumesToResponses(resumes *[]hh.ResumeDTO) (*[]Response, error) {
	responses := make([]Response, 0)
	for _, resume := range *resumes {
		for _, similarvac := range *resume.SimilarVacancies {
			responses = append(responses, Response{ResumeID: resume.ID, VacancyID: similarvac.ID, VacancyURL: similarvac.URL})
		}
	}
	if len(responses) == 0 {
		return nil, fmt.Errorf("no responses no resumes\n")
	}
	return &responses, nil
}
