package producer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/lookandhqte/workHelper/internal/entity"
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

func (p *TaskProducer) CreateTask(task *entity.Task) error {
	conn, err := beanstalk.Dial("tcp", p.addr)
	if err != nil {
		return fmt.Errorf("beanstalk dial failed: %w\n", err)
	}
	defer conn.Close()

	payload, err := json.Marshal(&task)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w\n", err)
	}

	_, err = conn.Put(payload, 1, 0, beforeNextTask*time.Second)

	return err
}
