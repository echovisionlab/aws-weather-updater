package main

import (
	"context"
	"github.com/echovisionlab/aws-weather-updater/pkg/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// prepare
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-exit
		cancel()
	}()
	app.Run(ctx)
}
