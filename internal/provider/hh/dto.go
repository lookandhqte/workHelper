package hh

//DTO структуры

// userInfoDTO структура информации о пользователе, возвращаемая запросом GET /me
type userInfoDTO struct {
	Email           string `json:"email"`
	ResumesURL      string `json:"resumes_url"`
	NegotiationsURL string `json:"negotiations_url"` //посмотреть GET количество откликов и приглашений
}

// userResumesDTO структура резюме пользователя, возвращаемая запросом GET ${ResumesURL}
type userResumesDTO struct {
	Items []resumeItem `json:"items"`
	Found int64        `json:"found"`
}

// resumeSimilarVacanciesDTO структура вакансий для пользователя, возвращаемая запросом GET ${Item.SimilarVacancies.URL}
type resumeSimilarVacanciesDTO struct {
	Items   []similarVacancyItem `json:"items"`
	Found   int64                `json:"found"`
	Pages   int64                `json:"pages"`
	Page    int64                `json:"page"`
	PerPage int64                `json:"per_page"`
}

// vacancyDataDTO структура необходимой информации для составления сопровода и отклика на вакансию
type vacancyDataDTO struct {
	ID                 string            `json:"id"`   //для отклика
	Name               string            `json:"name"` //искать отзывы (по желанию)
	SalaryRange        interface{}       `json:"salary_range"`
	Type               dictionary        `json:"type"`                  //должно быть open || anonymous
	Address            vacancyAddress    `json:"address"`               //адресс вакансии
	Experience         dictionary        `json:"experience"`            //под него ииха будет пиздеть
	Schedule           dictionary        `json:"schedule"`              //за счет него идет рассчет скок вилка на вакансию
	Employment         dictionary        `json:"employment"`            //исходя из этого будет считаться вилка
	Description        string            `json:"description"`           //для дипсика
	KeySkills          []vacancyKeySkill `json:"key_skills"`            //для дипсика
	ProfessionalRoles  []dictionary      `json:"professional_roles"`    //нужно 96 - можно все вакансии получать и отбирать по этому критерию
	Hidden             bool              `json:"hidden"`                //скрыта ли вакансия: должно быть false
	Employer           vacancyEmployer   `json:"employer"`              //инфа о работодателе (компании)
	Test               interface{}       `json:"test"`                  //тест
	AcceptTemporary    bool              `json:"accept_temporary"`      //доступность временного оформления: null или bool
	Approved           bool              `json:"approved"`              //прошла ли вакансия модерацию - должно быть true
	EmploymentForm     dictionary        `json:"employment_form"`       //все кроме project или fly_in_fly_out
	Internship         bool              `json:"internship"`            // стажа или нет
	WorkFormat         []dictionary      `json:"work_format"`           //интересует: гибрид удаленка офис
	WorkScheduleByDays []dictionary      `json:"work_schedule_by_days"` // 5/2 4/3 6/1 и тд
	WorkingHours       []dictionary      `json:"working_hours"`         // 6 7 8 в день
}

//Item структуры

// resumeItem структура описывающая экземпляр резюме
type resumeItem struct {
	URL              string           `json:"url"`               //URL самого резюме в hh
	ID               string           `json:"id"`                //id для отклика
	EmploymentForm   []dictionary     `json:"employment_form"`   // для иихи
	WorkFormat       []dictionary     `json:"work_format"`       //для иихи
	Education        userEducation    `json:"education"`         // для иихи
	Experience       []userExperience `json:"experience"`        // для иихи
	SimilarVacancies similarVacancies `json:"similar_vacancies"` // для поиска похожих вакансий
}

// similarVacancyItem структура описывающая экземпляр вакансии, получаемая запросом GET {$SimilarVacancies.URL}
type similarVacancyItem struct {
	ID       string `json:"id"`
	HasTest  bool   `json:"has_test"`
	Archived bool   `json:"archived"`
	URL      string `json:"url"`
}

// dictionary структура справочника
type dictionary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// userEducation структура описывающая образование соискателя
type userEducation struct {
	Level   dictionary    `json:"level"`
	Primary []userPrimary `json:"primary"`
}

// userPrimary структура описывающая единичное образование в определенной организации
type userPrimary struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Organization      string      `json:"organization"`
	Result            string      `json:"result"`
	Year              int64       `json:"year"`
	UniversityAcronym string      `json:"university_acronym"`
	NameID            string      `json:"name_id"`
	OrganizationID    interface{} `json:"organization_id"`
	ResultID          interface{} `json:"result_id"`
	EducationLevel    dictionary  `json:"education_level"`
}

// userExperience структура описывающая опыт соискателя
type userExperience struct {
	Start     string           `json:"start"`
	End       string           `json:"end"`
	Company   string           `json:"company"`
	CompanyID *string          `json:"company_id"`
	Employer  *vacancyEmployer `json:"employer"`
	Position  string           `json:"position"`
}

// similarVacancies вакансии, подходящие по резюме
type similarVacancies struct {
	URL      string   `json:"url"`
	Counters Counters `json:"counters"`
}

type Counters struct {
	Total int64 `json:"total"`
}

// vacancyAddress фактический адрес вакансии
type vacancyAddress struct {
	City string  `json:"city"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}

// vacancyEmployer структура информации о работодателе
type vacancyEmployer struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	URL                  string `json:"url"`
	AccreditedItEmployer bool   `json:"accredited_it_employer"`
	Trusted              bool   `json:"trusted"`
	Blacklisted          bool   `json:"blacklisted"`
}

// vacancyKeySkill структура ключевых скиллов вакансии
type vacancyKeySkill struct {
	Name string `json:"name"`
}

// Resume структура резюме
type ResumeDTO struct {
	ID               string                 `json:"id"`
	SimilarVacancies *[]SimilarVacanciesDTO `json:"similar_vacancies"`
}

type SimilarVacanciesDTO struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
