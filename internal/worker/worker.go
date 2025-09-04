package worker

import (
	"encoding/json"
	"log"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/lookandhqte/workHelper/internal/entity"
	accountUC "github.com/lookandhqte/workHelper/internal/usecase/account"
)

type UseCases struct {
	accoutUC accountUC.UseCase
}

type Worker struct {
	conn     *beanstalk.Conn
	stop     chan struct{}
	usecases UseCases
}

func NewWorker(addr string, uc accountUC.UseCase) *Worker {
	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Failed to connect to Beanstalkd:", err)
	}

	return &Worker{
		conn:     conn,
		stop:     make(chan struct{}),
		usecases: UseCases{accoutUC: uc},
	}
}

func (w *Worker) Start() {

	for {
		select {
		case <-w.stop:
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
	return nil
}

func (w *Worker) Stop() {
	close(w.stop)
	w.conn.Close()
}
