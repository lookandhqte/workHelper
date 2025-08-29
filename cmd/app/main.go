package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lookandhqte/workHelper/internal/app"
)

func main() {
	app := app.New()
	app.Run()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	app.Stop()
}
