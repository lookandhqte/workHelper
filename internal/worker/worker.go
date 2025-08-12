package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	id, body, err := conn.Reserve(timeForReserver * time.Second)
	if err != nil {
		return nil, nil
	}

	if err := conn.Delete(id); err != nil {
		return nil, fmt.Errorf("delete failed: %w", err)
	}

	return body, err
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: worker <beanstalk-addr>")
	}

	addr := os.Args[1]
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		<-sigChan

		cancel()

	}()

	worker := NewTaskWorker(addr)
	worker.Run(ctx)
}

func (w *TaskWorker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopped gracefully")
			return
		default:
			err := w.processTask(ctx)
			if err != nil {
				log.Printf("Error processing task: %v", err)
				select {
				case <-ctx.Done():
					return
				case <-time.After(1 * time.Second):
				}
			}
		}
	}
}

func (w *TaskWorker) processTask(ctx context.Context) error {
	conn, err := beanstalk.Dial("tcp", w.addr)
	if err != nil {
		return fmt.Errorf("beanstalk dial failed: %w", err)
	}

	defer conn.Close()

	id, _, err := conn.Reserve(5 * time.Second)
	if err != nil {
		if err == beanstalk.ErrTimeout {
			return nil
		}
		return fmt.Errorf("reserve failed: %v", err)
	}

	if err := conn.Delete(id); err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	return nil
}
