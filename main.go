package main

import (
	"github.com/echovisionlab/aws-weather-updater/pkg/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// prepare
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	app.Run(exit)
}
