package hh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/lookandhqte/workHelper/config"
	"github.com/lookandhqte/workHelper/internal/entity"
)

type Provider struct {
	app    *entity.App
	client *http.Client
}

// New создает нового провайдера
func New() *Provider {
	cfg := config.Load()
	return &Provider{
		app: &entity.App{
			RedirectURI:  cfg.RedirectURI,
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
		},
		client: &http.Client{},
	}
}

// GetToken меняет auth code на пару токенов, возвращает токены
func (r *Provider) GetToken(code string) (*entity.Token, error) {
	data := &url.Values{}
	data.Set("client_id", r.app.ClientID)
	data.Set("client_secret", r.app.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", r.app.RedirectURI)
	body, err := r.postRequest("https://api.hh.ru/token", data.Encode())
	if err != nil {
		fmt.Printf("err while rnew request func get token hh.go: %v\n", err)
		return nil, err
	}
	responseData := &entity.Token{}
	if err := json.Unmarshal(*body, responseData); err != nil {
		fmt.Printf("err while unmarshal body: %v\n", err)
		return nil, err
	}

	if responseData.AccessToken == "" && responseData.ExpiresIn == 0 && responseData.RefreshToken == "" {
		fmt.Printf("null result response data: %v\n", err)
		return nil, err
	}

	return responseData, nil
}

// RefreshToken обновляет токены
func (r *Provider) RefreshToken(refreshToken string) (*entity.Token, error) {
	data := &url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	body, err := r.postRequest("https://api.hh.ru/token", data.Encode())

	responseData := &entity.Token{}
	if err := json.Unmarshal(*body, responseData); err != nil {
		fmt.Printf("err while unmarshal body: %v\n", err)
		return nil, err
	}

	if responseData.AccessToken == "" && responseData.ExpiresIn == 0 && responseData.RefreshToken == "" {
		fmt.Printf("null result response data: %v\n", err)
		return nil, err
	}

	return responseData, nil
}

func (r *Provider) ReturnResponses(token string) (*[]Response, error) {

	userInfo, err := r.getUserInfo(token)
	if err != nil {
		fmt.Printf("err while rget user info: %v\n", err)
		return nil, err
	}

	resumesURL := userInfo.ResumesURL

	userResumes, err := r.getUserResumes(token, resumesURL)
	if err != nil {
		fmt.Printf("err while rget user info: %v\n", err)
		return nil, err
	}

	resumes := make([]ResumeDTO, 0, userResumes.Found)
	for _, resume := range userResumes.Items {
		if resume.ID == "ac4ac1e5ff0eefecc80039ed1f4f5555586a71" {
			found, err := r.getResumeSimilarVacancies(token, resume.SimilarVacancies.URL, resume.SimilarVacancies.Counters.Total)
			if err != nil {
				fmt.Printf("err while get list similar func fill account: %v\n", err)
				return nil, err
			}
			resumes = append(resumes, ResumeDTO{ID: resume.ID, SimilarVacancies: found})
		}
	}

	result, err := resumesToResponses(&resumes)
	if err != nil {
		fmt.Printf("err while get list similar func fill account: %v\n", err)
		return nil, err
	}

	return result, nil
}

func (r *Provider) ReturnPromptData(token string, responses *[]Response) (*[]PromptData, error) {
	result := make([]PromptData, 0, len(*responses))
	for _, response := range *responses {
		vacancy, err := r.getVacancy(token, response.VacancyURL)
		if err != nil {
			fmt.Printf("error while get vacancy return prompt data func: %v\n", err)
			return nil, err
		}
		prompt := preparePromptData(vacancy)
		result = append(result, *prompt)
	}
	return &result, nil
}

func (r *Provider) getRequest(token string, url string) (*[]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("err while rnew request func refresh token hh.go: %v\n", err)
		return nil, err
	}
	req.Header.Add("User-Agent", "workHelper/1.0(roselifemeow@gmail.com)")
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := r.client.Do(req)
	if err != nil {
		fmt.Printf("err while get user info dto: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err while readall auth get user info func %v\n", err)
		return nil, err
	}
	return &body, nil
}

func (r *Provider) postRequest(url string, data string) (*[]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data))
	if err != nil {
		fmt.Printf("err while rnew request func get token hh.go: %v\n", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data)))

	resp, err := r.client.Do(req)
	if err != nil {
		fmt.Printf("err while do req func get tokens hh.go: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err while readall get token func %v\n", err)
		return nil, err
	}
	return &body, nil
}

// getUserInfo получает информацию о пользователе
func (r *Provider) getUserInfo(token string) (*userInfoDTO, error) {
	body, err := r.getRequest(token, "https://api.hh.ru/me")
	if err != nil {
		fmt.Printf("err while readall auth get user info func %v\n", err)
		return nil, err
	}
	responseData := &userInfoDTO{}
	if err := json.Unmarshal(*body, responseData); err != nil {
		fmt.Printf("err while unmarshal body: %v\n", err)
		return nil, err
	}
	return responseData, nil
}

// GetUserResumes получает резюме пользователя
func (r *Provider) getUserResumes(token string, url string) (*userResumesDTO, error) {
	body, err := r.getRequest(token, url)
	if err != nil {
		fmt.Printf("err while readall auth get user info func %v\n", err)
		return nil, err
	}

	responseData := &userResumesDTO{}
	if err := json.Unmarshal(*body, responseData); err != nil {
		fmt.Printf("err while unmarshal body: %v\n", err)
		return nil, err
	}

	return responseData, nil
}

// getUserSimilarVacancies возвращает id и url всех подходящих к резюме вакансий
func (r *Provider) getResumeSimilarVacancies(token string, url string, amountOfVacancies int64) (*[]SimilarVacanciesDTO, error) {
	similarVacancies := make([]SimilarVacanciesDTO, 0, amountOfVacancies)
	var page int64 = 0
	for {
		paginatedURL := fmt.Sprintf("%s?page=%d&per_page=50", url, page)

		body, err := r.getRequest(token, paginatedURL)
		if err != nil {
			fmt.Printf("err while get request: %v\n", err)
			return nil, err
		}

		responseData := &resumeSimilarVacanciesDTO{}
		if err := json.Unmarshal(*body, responseData); err != nil {
			fmt.Printf("err while unmarshal body: %v\n", err)
			return nil, err
		}

		for _, item := range responseData.Items {
			if !item.Archived {
				similarVacancies = append(similarVacancies, SimilarVacanciesDTO{ID: item.ID, URL: item.URL})
			}
		}

		if len(similarVacancies) >= int(amountOfVacancies) ||
			page >= responseData.Pages-1 ||
			len(responseData.Items) == 0 {
			break
		}

		page++
	}

	return &similarVacancies, nil
}

func (r *Provider) getVacancy(token string, url string) (*vacancyDataDTO, error) {
	body, err := r.getRequest(token, url)
	if err != nil {
		fmt.Printf("err while readall auth get user info func %v\n", err)
		return nil, err
	}

	responseData := &vacancyDataDTO{}
	if err := json.Unmarshal(*body, responseData); err != nil {
		fmt.Printf("err while unmarshal body: %v\n", err)
		return nil, err
	}

	return responseData, nil
}

func (r *Provider) DoResponse(token string, responseData *Response) error {

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("message", responseData.Message); err != nil {
		return err
	}
	if err := writer.WriteField("resume_id", responseData.ResumeID); err != nil {
		return err
	}
	if err := writer.WriteField("vacancy_id", responseData.VacancyID); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.hh.ru/negotiations", &body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "workHelper/1.0 (roselifemeow@gmail.com)")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	fmt.Printf("Success response: %s\n", string(respBody))
	return nil
}
