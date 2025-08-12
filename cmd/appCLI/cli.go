package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/config"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/app"
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

					beanstalkAddr := cfg.BeanstalkAddr
					processes := make([]*exec.Cmd, 0, numWorkers)

					for i := 0; i < numWorkers; i++ {
						cmd := exec.Command("./worker", beanstalkAddr)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr

						if err := cmd.Start(); err != nil {
							log.Printf("Failed to start worker %d: %v", i+1, err)
							continue
						}

						processes = append(processes, cmd)
						log.Printf("Started worker %d (PID: %d)", i+1, cmd.Process.Pid)
					}

					sigChan := make(chan os.Signal, 1)
					signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
					<-sigChan

					log.Println("Stopping workers...")
					for _, cmd := range processes {
						if cmd.Process != nil {
							cmd.Process.Signal(syscall.SIGTERM)
						}
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
