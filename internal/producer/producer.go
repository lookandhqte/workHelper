package producer

import (
	"encoding/json"
	"fmt"
	"time"

	v1 "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/controller/http/v1"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"github.com/beanstalkd/go-beanstalk"
)

type TaskProducer struct {
	addr string
}

func NewTaskProducer(addr string) *TaskProducer {
	return &TaskProducer{addr: addr}
}

const (
	beforeNextTask = 10
)

func (p *TaskProducer) EnqueueSyncContactsTask(accountID int, integrationID int, contacts []entity.GlobalContact) error {
	conn, err := beanstalk.Dial("tcp", p.addr)
	if err != nil {
		return fmt.Errorf("beanstalk dial failed: %w", err)
	}
	defer conn.Close()

	task := struct {
		AccountID     int                    `json:"account_id"`
		IntegrationID int                    `json:"integration_id"`
		Contacts      []entity.GlobalContact `json:"contacts"`
	}{
		AccountID:     accountID,
		IntegrationID: integrationID,
		Contacts:      contacts,
	}

	payload, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	_, err = conn.Put(payload, 1, 0, beforeNextTask*time.Second)
	return err
}

func (p *TaskProducer) EnqueueSyncWebhookContactsTask(webhookData v1.WebhookContactDTO) error {
	conn, err := beanstalk.Dial("tcp", p.addr)
	if err != nil {
		return fmt.Errorf("beanstalk dial failed: %w", err)
	}
	defer conn.Close()

	tasks := make([]struct {
		AccountID string `json:"account_id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}, 0, len(webhookData.Contacts.Add))

	for _, apiContact := range webhookData.Contacts.Add {
		task := struct {
			AccountID string `json:"account_id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Phone     string `json:"phone"`
		}{
			AccountID: apiContact.AccountID,
			Name:      apiContact.Name,
		}

		for _, field := range apiContact.CustomFields {
			if field.Code == "EMAIL" && len(field.Values) > 0 {
				task.Email = field.Values[0].Value
			}
		}

		tasks = append(tasks, task)
	}

	payload, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	_, err = conn.Put(payload, 1, 0, 120*time.Second)
	return err
}
