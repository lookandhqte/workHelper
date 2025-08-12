package producer

import (
	"encoding/json"
	"fmt"
	"time"

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
	beforeNextTask = 120
)

func (p *TaskProducer) EnqueueCreateContactTask(globalContact *entity.GlobalContact) error {
	conn, err := beanstalk.Dial("tcp", p.addr)
	if err != nil {
		return fmt.Errorf("beanstalk dial failed: %w", err)
	}
	defer conn.Close()

	payload, err := json.Marshal(globalContact)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	_, err = conn.Put(payload, 1, 0, beforeNextTask*time.Second)
	return err
}
