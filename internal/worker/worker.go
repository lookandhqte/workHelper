package worker

import (
	"fmt"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

type TaskWorker struct {
	addr string
}

func NewTaskWorker(addr string) *TaskWorker {
	return &TaskWorker{addr: addr}
}

const (
	timeForReserver = 5
)

func (p *TaskWorker) ResolveCreateContactTask() ([]byte, error) {
	conn, err := beanstalk.Dial("tcp", p.addr)
	if err != nil {
		return nil, fmt.Errorf("beanstalk dial failed: %w", err)
	}
	defer conn.Close()

	_, body, err := conn.Reserve(timeForReserver * time.Second)
	return body, err
}
