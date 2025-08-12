package worker

import (
	"fmt"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

type TaskWorker struct {
	addr   string
	status string
}

type TaskWorkers []TaskWorker

func NewTaskWorkers(addr string, amount int) *TaskWorkers {
	workers := make(TaskWorkers, 0, amount)
	for len(workers) < amount {
		workers = append(workers, TaskWorker{addr: addr, status: "available"})
	}
	return &workers
}

const (
	timeForReserver = 5
	defaultDuration = 10 * time.Second
)

func (p *TaskWorker) ResolveCreateContactTask() ([]byte, error) {
	p.status = "busy"
	conn, err := beanstalk.Dial("tcp", p.addr)
	if err != nil {
		return nil, fmt.Errorf("beanstalk dial failed: %w", err)
	}
	defer conn.Close()

	_, body, err := conn.Reserve(timeForReserver * time.Second)
	p.status = "available"
	return body, err
}

func (p *TaskWorkers) GetAvailableWorker(workers *TaskWorkers) *TaskWorker {
	return p.ReturnWorkerWithTimeout(workers, defaultDuration)
}

func (p *TaskWorkers) ReturnWorkerWithTimeout(workers *TaskWorkers, timeout time.Duration) *TaskWorker {
	for {
		for _, worker := range *workers {
			if worker.status == "available" {
				return &worker
			}
		}
		time.Sleep(timeout)
	}
}
