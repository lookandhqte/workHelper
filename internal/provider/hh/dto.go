package hh

// userInfoDTO структура информации о пользователе, возвращаемая запросом GET /me
type userInfoDTO struct {
	Email           string `json:"email"`
	ResumesURL      string `json:"resumes_url"`
	NegotiationsURL string `json:"negotiations_url"`
}

// userResumesDTO структура резюме пользователя, возвращаемая запросом GET ${ResumesURL}
type userResumesDTO struct {
	Items []resumeItem `json:"items"`
	Found int64        `json:"found"`
}

// userSimilarVacancyDTO структура вакансий для пользователя, возвращаемая запросом GET ${Item.SimilarVacancies.URL}
type userSimilarVacancyDTO struct {
	Items []similarVacancyItem `json:"items"`
	Found int64                `json:"found"`
}

// resumeItem структура описывающая экземпляр резюме
type resumeItem struct {
	URL              string           `json:"url"` //URL самого резюме в hh
	ID               string           `json:"id"`
	EmploymentForm   []Gender         `json:"employment_form"`   // для иихи
	WorkFormat       []Gender         `json:"work_format"`       //для иихи
	Education        Education        `json:"education"`         // для иихи
	Experience       []Experience     `json:"experience"`        // для иихи
	SimilarVacancies SimilarVacancies `json:"similar_vacancies"` // для поиска похожих вакансий
}

// vacancyItem структура описывающая экземпляр вакансии, получаемая запросом GET {$SimilarVacancies.URL}
type similarVacancyItem struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	HasTest                bool   `json:"has_test"`
	ResponseLetterRequired bool   `json:"response_letter_required"`
	Archived               bool   `json:"archived"`
	URL                    string `json:"url"`
	NightShifts            bool   `json:"night_shifts"` //это важно вроде
	Internship             bool   `json:"internship"`
}

// Gender тоже хз для чего
type Gender struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Area тоже хз для чего
type Area struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Education структура описывающая образование соискателя
type Education struct {
	Level   Gender    `json:"level"`
	Primary []Primary `json:"primary"`
}

// Primary структура описывающая единичное образование в определенной организации
type Primary struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Organization      string      `json:"organization"`
	Result            string      `json:"result"`
	Year              int64       `json:"year"`
	UniversityAcronym string      `json:"university_acronym"`
	NameID            string      `json:"name_id"`
	OrganizationID    interface{} `json:"organization_id"`
	ResultID          interface{} `json:"result_id"`
	EducationLevel    Gender      `json:"education_level"`
}

// Experience структура описывающая опыт соискателя
type Experience struct {
	Start      string      `json:"start"`
	End        string      `json:"end"`
	Company    string      `json:"company"`
	CompanyID  *string     `json:"company_id"`
	Industry   interface{} `json:"industry"`
	Industries []Gender    `json:"industries"`
	Area       interface{} `json:"area"`
	CompanyURL interface{} `json:"company_url"`
	Employer   *Employer   `json:"employer"`
	Position   string      `json:"position"`
}

// Employer - в будущем по имени нужно будет искать отзывы и сведения для составления образа и релевантности компании
type Employer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type SimilarVacancies struct {
	URL string `json:"url"`
}
