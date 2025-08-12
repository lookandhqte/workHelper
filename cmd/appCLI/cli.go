package main

import (
	"context"
	"log"
	"os"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/app"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/worker"
	"github.com/urfave/cli/v2"
)

func main() {
	ourApp := app.New()
	cfg := config.Load()

	app := &cli.App{
		Name:  "amocrm-app",
		Usage: "AMOCRM integration service",
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "Start HTTP server",
				Action: func(c *cli.Context) error {
					ourApp.Run()
					return nil
				},
			},
			{
				Name: "worker",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "workers",
						Aliases:  []string{"n"},
						Usage:    "Number of workers to start",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					numWorkers := c.Int("workers")
					ourApp.AddTask(func(ctx context.Context) {
						workers := worker.NewTaskWorkers(cfg.BeanstalkAddr, numWorkers)

						for i := 0; i < numWorkers; i++ {
							go func() {
								defer func() {
									if r := recover(); r != nil {
										log.Printf("Worker panic: %v", r)
									}
								}()
								for {
									select {
									case <-ctx.Done():
										return
									default:
										_, err := workers.ResolveCreateContactTask(workers)
										if err != nil {
											log.Printf("Worker error: %v", err)
										}
									}
								}
							}()
						}

						<-ctx.Done()
						log.Println("All workers stopped")
					})
					return nil

				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
