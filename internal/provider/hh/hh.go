package hh

import (
	"encoding/json"
	"fmt"
	"io"
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

	req, err := http.NewRequest(http.MethodPost, "https://api.hh.ru/token", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Printf("err while rnew request func get token hh.go: %v\n", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

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

	responseData := &entity.Token{}
	if err := json.Unmarshal(body, responseData); err != nil {
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

	req, err := http.NewRequest(http.MethodPost, "https://api.hh.ru/token", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Printf("err while rnew request func refresh token hh.go: %v\n", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := r.client.Do(req)
	if err != nil {
		fmt.Printf("err while do req func refresh tokens hh.go: %v\n", err)
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

	responseData := &entity.Token{}
	if err := json.Unmarshal(body, responseData); err != nil {
		fmt.Printf("err while unmarshal body: %v\n", err)
		return nil, err
	}

	if responseData.AccessToken == "" && responseData.ExpiresIn == 0 && responseData.RefreshToken == "" {
		fmt.Printf("null result response data: %v\n", err)
		return nil, err
	}

	return responseData, nil
}

// getUserInfo получает информацию о пользователе
func (r *Provider) getUserInfo(token string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.hh.ru/me", nil)
	if err != nil {
		fmt.Printf("err while rnew request func refresh token hh.go: %v\n", err)
		return "", err
	}
	req.Header.Add("User-Agent", "workHelper/1.0(roselifemeow@gmail.com)")
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := r.client.Do(req)
	if err != nil {
		fmt.Printf("err while get user info dto: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err while readall auth get user info func %v\n", err)
		return "", err
	}

	responseData := &userInfoDTO{}
	if err := json.Unmarshal(body, responseData); err != nil {
		fmt.Printf("err while unmarshal body: %v\n", err)
		return "", err
	}
	fmt.Printf("response after unmarshal %v\n", responseData)
	//fmt.Println(responseData.Email)
	return responseData.ResumesURL, nil
}

// GetUserResumes получает информацию о пользователе
func (r *Provider) getUserResumes(token string) (*[]string, error) {
	baseURL, err := r.getUserInfo(token)
	if err != nil {
		fmt.Printf("err while get user info: %v\n", err)
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
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

	responseData := &userResumesDTO{}
	if err := json.Unmarshal(body, responseData); err != nil {
		fmt.Printf("err while unmarshal body: %v\n", err)
		return nil, err
	}
	fmt.Println("found resumes:")
	fmt.Println(responseData.Found)
	similarVacanciesURLs := make([]string, 0, responseData.Found)
	for _, item := range responseData.Items {
		similarVacanciesURLs = append(similarVacanciesURLs, item.SimilarVacancies.URL)
	}
	return &similarVacanciesURLs, nil
}

// getUserSimilarVacanciesIDs получает информацию о пользователе
func (r *Provider) getUserSimilarVacanciesURLs(token string) (*[]string, error) {
	baseURL, err := r.getUserResumes(token)
	foundURLs := make([]string, 0, len(*baseURL))
	if err != nil {
		fmt.Printf("err while get user resumes: %v\n", err)
		return nil, err
	}
	for _, url := range *baseURL {
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
		responseData := &userSimilarVacancyDTO{}
		if err := json.Unmarshal(body, responseData); err != nil {
			fmt.Printf("err while unmarshal body: %v\n", err)
			return nil, err
		}
		fmt.Println("found vacancies:")
		fmt.Println(responseData.Found)
		for _, item := range responseData.Items {
			if !item.Archived && !item.HasTest {
				foundURLs = append(foundURLs, item.URL)
			}
		}
	}

	return &foundURLs, nil
}

// // getUserVacanciesIDs возвращает вакансии отобранные по своим алгоритмам
// func (r *Provider) getUserVacanciesURLs(token string) (*[]string, error) {

// 	data := &url.Values{}
// 	data.Set("professional_roles", "96")

// 		req, err := http.NewRequest(http.MethodGet, "https://api.hh.ru/vacancies", strings.NewReader(data.Encode()))
// 		if err != nil {
// 			fmt.Printf("err while rnew request func refresh token hh.go: %v\n", err)
// 			return nil, err
// 		}
// 		req.Header.Add("User-Agent", "workHelper/1.0(roselifemeow@gmail.com)")
// 		req.Header.Add("Authorization", "Bearer "+token)
// 		resp, err := r.client.Do(req)
// 		if err != nil {
// 			fmt.Printf("err while get user info dto: %v\n", err)
// 			return nil, err
// 		}
// 		defer resp.Body.Close()

// 		if resp.StatusCode != http.StatusOK {
// 			body, _ := io.ReadAll(resp.Body)
// 			return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
// 		}
// 		body, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Printf("err while readall auth get user info func %v\n", err)
// 			return nil, err
// 		}
// 		responseData := &userSimilarVacancyDTO{}
// 		if err := json.Unmarshal(body, responseData); err != nil {
// 			fmt.Printf("err while unmarshal body: %v\n", err)
// 			return nil, err
// 		}
// 		fmt.Println("found vacancies:")
// 		fmt.Println(responseData.Found)

// 	return nil, nil
// }

// getVacancyData возвращает описание вакансии по ID
func (r *Provider) GetVacancyData(token string) (*[]vacancyDataDTO, error) {
	vacancyURLs, err := r.getUserSimilarVacanciesURLs(token)
	if err != nil {
		fmt.Printf("err while get user similar vacancies urls: %v\n", err)
		return nil, err
	}

	vacanciesData := make([]vacancyDataDTO, 0, len(*vacancyURLs))

	for _, url := range *vacancyURLs {
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
		responseData := &vacancyDataDTO{}
		if err := json.Unmarshal(body, responseData); err != nil {
			fmt.Printf("err while unmarshal body: %v\n", err)
			return nil, err
		}

		vacanciesData = append(vacanciesData, *responseData)
	}
	return &vacanciesData, nil
}
