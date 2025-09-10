package deepseek

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lookandhqte/workHelper/config"
	"github.com/lookandhqte/workHelper/internal/provider/hh"
)

type Provider struct {
	key    string
	client *http.Client
}

// New создает нового провайдера
func New() *Provider {
	cfg := config.Load()
	return &Provider{
		key:    cfg.DeepseekAPI,
		client: &http.Client{},
	}
}

// GetVacancySoprovod возвращает сопроводительное для вакансии
func (r *Provider) GetVacancySoprovod(promptData hh.PromptData) (string, error) {
	data := fmt.Sprintf(`{
    "model": "deepseek-chat",
    "messages": [
        {"role": "system", "content": "You are a helpful assistant. You specialize on writing cover letters based on the job description. Your style is informative and concise, without unnecessary words. You use russian language for writing cover letters."},
        {"role": "user", "content": "Напиши короткое сопроводительное письмо для отклика на вакансию. Пиши в женском роде, я женщина. ДАННЫЕ ВАКАНСИИ: Должность: %s Компания: %s Требуемый опыт: %s Ключевые навыки, требуемые в вакансии: %v Описание вакансии: %s ТРЕБОВАНИЯ К ПИСЬМУ: 1. Максимально коротко и по делу (3-4 предложения) 2. Подчеркнуть соответствие ключевым навыкам 3. Указать заинтересованность в позиции 4. Профессиональный тон без излишней эмоциональности 5. В конце добавить: 'Этот отклик был отправлен автоматически, разработанной мной CRM-системой, на данный момент я тестирую отправку откликов на вакансии и автогенерацию сопроводительных. Спасибо!'"}
    ],
    "stream": false
}`, promptData.Position, promptData.Company, promptData.Experience, promptData.KeySkills, promptData.Description)

	req, err := http.NewRequest(http.MethodPost, "https://api.deepseek.com/chat/completions", strings.NewReader(data))
	if err != nil {
		fmt.Printf("err while rnew request deepseek: %v\n", err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+r.key)
	resp, err := r.client.Do(req)
	if err != nil {
		fmt.Printf("err while do req func refresh tokens hh.go: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err while readall get token func %v\n", err)
		return "", err
	}
	deepseekResp := &responseDTO{}
	if err := json.Unmarshal(body, deepseekResp); err != nil {
		fmt.Printf("err while unmarshal body: %v\n", err)
		return "", err
	}
	return deepseekResp.Choices[0].Message.Content, nil
}
