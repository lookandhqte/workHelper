package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/lookandhqte/workHelper/internal/entity"
	"github.com/lookandhqte/workHelper/internal/provider"
	accountUC "github.com/lookandhqte/workHelper/internal/usecase/account"
	tokenUC "github.com/lookandhqte/workHelper/internal/usecase/token"
)

type UseCases struct {
	accoutUC accountUC.UseCase
	tokenUC  tokenUC.UseCase
}

type Worker struct {
	conn     *beanstalk.Conn
	stop     chan struct{}
	usecases UseCases
	provider provider.Provider
}

func NewWorker(addr string, uc accountUC.UseCase, tuc tokenUC.UseCase, provider provider.Provider) *Worker {
	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Failed to connect to Beanstalkd:", err)
	}

	return &Worker{
		conn:     conn,
		stop:     make(chan struct{}),
		usecases: UseCases{accoutUC: uc, tokenUC: tuc},
		provider: provider,
	}
}

func (w *Worker) Start() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		case <-w.stop:
			cancel()
			return
		default:
			id, body, err := w.conn.Reserve(5 * time.Second)
			if err != nil {
				continue
			}

			if err := w.processTask(body); err != nil {
				log.Printf("Task failed: %v", err)
				w.conn.Release(id, 0, 0)
				continue
			}

			w.conn.Delete(id)
		}
	}
}

func (w *Worker) processTask(body []byte) error {
	var task entity.Task
	if err := json.Unmarshal(body, &task); err != nil {
		return err
	}

	switch task.Type {
	case "account_creating":
		return w.handleAccountCreation(task.Payload)
	case "account_updating":
		return w.handleAccountUpdating(task.Payload)
	}

	return nil
}

func (w *Worker) handleAccountCreation(body []byte) error {
	var account entity.Account
	if err := json.Unmarshal(body, &account); err != nil {
		log.Printf("er while unmarshal: %v\n", err)
		return err
	}

	if err := w.usecases.accoutUC.Create(&account); err != nil {
		log.Printf("er while creting func handle: %v\n", err)
		return err
	}
	return nil
}

func (w *Worker) handleAccountUpdating(body []byte) error {
	var account entity.Account
	if err := json.Unmarshal(body, &account); err != nil {
		log.Printf("er while unmarshal: %v\n", err)
		return err
	}

	if err := w.usecases.accoutUC.Update(&account); err != nil {
		log.Printf("er while creting func handle update: %v\n", err)
		return err
	}

	go w.scheduleTokenRefresh()
	return nil
}

func (w *Worker) scheduleTokenRefresh() {
	account, err := w.usecases.accoutUC.Return()
	if err != nil {
		fmt.Printf("lox ne poluchilos: %v\n", err)
	}
	if account.Token.ExpiresIn != 0 {
		expiresAt := time.Unix(int64(account.Token.CreatedAt), 0).Add(
			time.Duration(account.Token.ExpiresIn) * time.Second,
		)
		refreshTime := expiresAt.Add(-1 * time.Hour)
		<-time.After(time.Until(refreshTime))

		if err := w.refreshToken(); err != nil {
			log.Printf("failed to refresh token: %v", err)
		}
	}
}

func (w *Worker) refreshToken() error {
	account, err := w.usecases.accoutUC.Return()
	if err != nil {
		return err
	}
	newTokens, err := w.provider.HH.RefreshToken(account.Token.RefreshToken)
	if err != nil {
		return err
	}
	newTokens.CreatedAt = int(time.Now().Unix())
	newTokens.AccountID = account.ID
	w.usecases.tokenUC.Create(newTokens)
	account.Token = *newTokens
	return w.usecases.accoutUC.Update(account)
}

func (w *Worker) Stop() {
	close(w.stop)
	w.conn.Close()
}
