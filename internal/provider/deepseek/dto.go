package deepseek

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

// dictionary структура справочника
type dictionary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

type responseDTO struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
