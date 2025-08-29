package main

import (
	"os"
	"os/signal"
	"syscall"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/app"
)

func main() {
	app := app.New()
	app.Run()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	app.Stop()
}
