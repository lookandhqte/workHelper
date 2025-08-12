package worker

import (
	"fmt"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

type TaskWorkers []struct {
	addr   string
	status string
}

func NewTaskWorkers(addr string, amount int) *TaskWorkers {
	workers := make(TaskWorkers, 0, amount)
	for len(workers) < amount {
		workers = append(workers, struct {
			addr   string
			status string
		}{addr: addr, status: "available"})
	}
	return &workers
}

const (
	timeForReserver = 5
	defaultDuration = 10 * time.Second
)

func (p *TaskWorkers) ResolveCreateContactTask(workers *TaskWorkers) ([]byte, error) {
	worker := p.GetAvailableWorker(workers)
	worker.status = "busy"
	defer func() { worker.status = "available" }()
	conn, err := beanstalk.Dial("tcp", worker.addr)
	if err != nil {
		return nil, fmt.Errorf("beanstalk dial failed: %w", err)
	}
	defer conn.Close()

	_, body, err := conn.Reserve(timeForReserver * time.Second)
	return body, err
}

func (p *TaskWorkers) GetAvailableWorker(workers *TaskWorkers) *struct {
	addr   string
	status string
} {
	return p.ReturnWorkerWithTimeout(workers, defaultDuration)
}

func (p *TaskWorkers) ReturnWorkerWithTimeout(workers *TaskWorkers, timeout time.Duration) *struct {
	addr   string
	status string
} {
	for {
		for _, worker := range *workers {
			if worker.status == "available" {
				return &worker
			}
		}
		time.Sleep(timeout)
	}
}
