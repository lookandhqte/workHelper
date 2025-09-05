package hh

type userInfoDTO struct {
	AuthType              string            `json:"auth_type"`
	IsApplicant           bool              `json:"is_applicant"`
	IsEmployer            bool              `json:"is_employer"`
	IsAdmin               bool              `json:"is_admin"`
	IsHiringManager       bool              `json:"is_hiring_manager"`
	IsApplication         bool              `json:"is_application"`
	IsEmployerIntegration bool              `json:"is_employer_integration"`
	CryptedID             string            `json:"crypted_id"`
	ID                    string            `json:"id"`
	IsAnonymous           bool              `json:"is_anonymous"`
	Email                 string            `json:"email"`
	FirstName             string            `json:"first_name"`
	MiddleName            string            `json:"middle_name"`
	LastName              string            `json:"last_name"`
	ResumesURL            string            `json:"resumes_url"`
	NegotiationsURL       string            `json:"negotiations_url"`
	IsInSearch            bool              `json:"is_in_search"`
	MidName               string            `json:"mid_name"`
	Employer              interface{}       `json:"employer"`
	Manager               interface{}       `json:"manager"`
	Phone                 string            `json:"phone"`
	Counters              userCounters      `json:"counters"`
	ProfileVideos         userProfileVideos `json:"profile_videos"`
	PersonalManager       interface{}       `json:"personal_manager"`
}

type userCounters struct {
	NewResumeViews     int64 `json:"new_resume_views"`
	UnreadNegotiations int64 `json:"unread_negotiations"`
	ResumesCount       int64 `json:"resumes_count"`
}

type userProfileVideos struct {
	Items []interface{} `json:"items"`
}

// Резюме
type userResumesDTO struct {
	Items []resumeItem `json:"items"`
	Found int64        `json:"found"`
}

type resumeItem struct {
	URL              string           `json:"url"`
	AlternateURL     string           `json:"alternate_url"`
	ID               string           `json:"id"`
	Platform         Platform         `json:"platform"`
	EmploymentForm   []Gender         `json:"employment_form"`
	WorkFormat       []Gender         `json:"work_format"`
	RealID           string           `json:"real_id"`
	Education        Education        `json:"education"`
	Experience       []Experience     `json:"experience"`
	Marked           bool             `json:"marked"`
	Finished         bool             `json:"finished"`
	Status           Gender           `json:"status"`
	Access           Access           `json:"access"`
	NextPublishAt    string           `json:"next_publish_at"`
	Contact          []Contact        `json:"contact"`
	Tags             []interface{}    `json:"tags"`
	Visible          bool             `json:"visible"`
	Created          string           `json:"created"`
	Updated          string           `json:"updated"`
	SimilarVacancies SimilarVacancies `json:"similar_vacancies"`
	NewViews         int64            `json:"new_views"`
	TotalViews       int64            `json:"total_views"`
	ViewsURL         string           `json:"views_url"`
}

type Access struct {
	Type Gender `json:"type"`
}

type Gender struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Actions struct {
	Download Download `json:"download"`
}

type Download struct {
	PDF PDF `json:"pdf"`
	Rtf PDF `json:"rtf"`
}

type PDF struct {
	URL string `json:"url"`
}

type Area struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Contact struct {
	Value            interface{} `json:"value"`
	Type             Gender      `json:"type"`
	Preferred        bool        `json:"preferred"`
	Comment          *string     `json:"comment"`
	ContactValue     string      `json:"contact_value"`
	Kind             string      `json:"kind"`
	NeedVerification *bool       `json:"need_verification,omitempty"`
	Verified         *bool       `json:"verified,omitempty"`
}

type ValueClass struct {
	Country   string `json:"country"`
	City      string `json:"city"`
	Number    string `json:"number"`
	Formatted string `json:"formatted"`
}

type Education struct {
	Level   Gender    `json:"level"`
	Primary []Primary `json:"primary"`
}

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

type Employer struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	URL          string   `json:"url"`
	AlternateURL string   `json:"alternate_url"`
	LogoUrls     LogoUrls `json:"logo_urls"`
}

type PaidService struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type Photo struct {
	ID     string `json:"id"`
	Small  string `json:"small"`
	Medium string `json:"medium"`
	The40  string `json:"40"`
	The100 string `json:"100"`
	The500 string `json:"500"`
}

type Platform struct {
	ID string `json:"id"`
}

type SimilarVacancies struct {
	URL      string   `json:"url"`
	Counters Counters `json:"counters"`
}

type Counters struct {
	Total int64 `json:"total"`
}

type TotalExperience struct {
	Months int64 `json:"months"`
}

type ValueUnion struct {
	String     *string
	ValueClass *ValueClass
}

// Вакансии
type vacancyDTO struct {
	Items        []vacancyItem `json:"items"`
	Found        int64         `json:"found"`
	Pages        int64         `json:"pages"`
	Page         int64         `json:"page"`
	PerPage      int64         `json:"per_page"`
	Clusters     interface{}   `json:"clusters"`
	Arguments    interface{}   `json:"arguments"`
	Fixes        interface{}   `json:"fixes"`
	Suggests     interface{}   `json:"suggests"`
	AlternateURL string        `json:"alternate_url"`
}

type vacancyItem struct {
	ID                      string              `json:"id"`
	Premium                 bool                `json:"premium"`
	Name                    string              `json:"name"`
	Department              interface{}         `json:"department"`
	HasTest                 bool                `json:"has_test"`
	ResponseLetterRequired  bool                `json:"response_letter_required"`
	Area                    Area                `json:"area"`
	Salary                  *vacancySalary      `json:"salary"`
	SalaryRange             *vacancySalaryRange `json:"salary_range"`
	Type                    vacancyEmployment   `json:"type"`
	Address                 *vacancyAddress     `json:"address"`
	ResponseURL             interface{}         `json:"response_url"`
	SortPointDistance       interface{}         `json:"sort_point_distance"`
	PublishedAt             string              `json:"published_at"`
	CreatedAt               string              `json:"created_at"`
	Archived                bool                `json:"archived"`
	ApplyAlternateURL       string              `json:"apply_alternate_url"`
	ShowLogoInSearch        *bool               `json:"show_logo_in_search"`
	ShowContacts            bool                `json:"show_contacts"`
	InsiderInterview        interface{}         `json:"insider_interview"`
	URL                     string              `json:"url"`
	AlternateURL            string              `json:"alternate_url"`
	Relations               []interface{}       `json:"relations"`
	Employer                vacancyEmployer     `json:"employer"`
	Snippet                 vacancySnippet      `json:"snippet"`
	Contacts                *vacancyContacts    `json:"contacts"`
	Schedule                vacancyEmployment   `json:"schedule"`
	WorkingDays             []vacancyEmployment `json:"working_days"`
	WorkingTimeIntervals    []vacancyEmployment `json:"working_time_intervals"`
	WorkingTimeModes        []interface{}       `json:"working_time_modes"`
	AcceptTemporary         bool                `json:"accept_temporary"`
	FlyInFlyOutDuration     []interface{}       `json:"fly_in_fly_out_duration"`
	WorkFormat              []vacancyEmployment `json:"work_format"`
	WorkingHours            []vacancyEmployment `json:"working_hours"`
	WorkScheduleByDays      []vacancyEmployment `json:"work_schedule_by_days"`
	NightShifts             bool                `json:"night_shifts"`
	ProfessionalRoles       []vacancyEmployment `json:"professional_roles"`
	AcceptIncompleteResumes bool                `json:"accept_incomplete_resumes"`
	Experience              vacancyEmployment   `json:"experience"`
	Employment              vacancyEmployment   `json:"employment"`
	EmploymentForm          vacancyEmployment   `json:"employment_form"`
	Internship              bool                `json:"internship"`
	AdvResponseURL          interface{}         `json:"adv_response_url"`
	IsAdvVacancy            bool                `json:"is_adv_vacancy"`
	AdvContext              interface{}         `json:"adv_context"`
	Branding                *vacancyBranding    `json:"branding,omitempty"`
}

type vacancyAddress struct {
	City          string         `json:"city"`
	Street        *string        `json:"street"`
	Building      *string        `json:"building"`
	Lat           float64        `json:"lat"`
	Lng           float64        `json:"lng"`
	Description   interface{}    `json:"description"`
	Raw           string         `json:"raw"`
	Metro         *vacancyMetro  `json:"metro"`
	MetroStations []vacancyMetro `json:"metro_stations"`
	ID            string         `json:"id"`
}

type vacancyMetro struct {
	StationName string  `json:"station_name"`
	LineName    string  `json:"line_name"`
	StationID   string  `json:"station_id"`
	LineID      string  `json:"line_id"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
}

type vacancyBranding struct {
	Type   string      `json:"type"`
	Tariff interface{} `json:"tariff"`
}

type vacancyContacts struct {
	Name                string         `json:"name"`
	Email               string         `json:"email"`
	Phones              []vacancyPhone `json:"phones"`
	CallTrackingEnabled bool           `json:"call_tracking_enabled"`
}

type vacancyPhone struct {
	Comment   *string `json:"comment"`
	City      string  `json:"city"`
	Number    string  `json:"number"`
	Country   string  `json:"country"`
	Formatted string  `json:"formatted"`
}

type vacancyEmployer struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	URL                  string    `json:"url"`
	AlternateURL         string    `json:"alternate_url"`
	LogoUrls             *LogoUrls `json:"logo_urls"`
	VacanciesURL         string    `json:"vacancies_url"`
	AccreditedItEmployer bool      `json:"accredited_it_employer"`
	Trusted              bool      `json:"trusted"`
}

type vacancyEmployment struct {
	ID   string      `json:"id"`
	Name vacancyName `json:"name"`
}

type vacancySalary struct {
	From     *int64 `json:"from"`
	To       *int64 `json:"to"`
	Currency string `json:"currency"`
	Gross    bool   `json:"gross"`
}

type vacancySalaryRange struct {
	From      *int64             `json:"from"`
	To        *int64             `json:"to"`
	Currency  string             `json:"currency"`
	Gross     bool               `json:"gross"`
	Mode      vacancyEmployment  `json:"mode"`
	Frequency *vacancyEmployment `json:"frequency"`
}

type vacancySnippet struct {
	Requirement    string  `json:"requirement"`
	Responsibility *string `json:"responsibility"`
}

type vacancyName string

type LogoUrls struct {
	The90    string `json:"90"`
	The240   string `json:"240"`
	Original string `json:"original"`
}
