package app

import (
	"context"
	"fmt"
	"github.com/echovisionlab/aws-weather-updater/pkg/browser"
	"github.com/echovisionlab/aws-weather-updater/pkg/database"
	"github.com/echovisionlab/aws-weather-updater/pkg/task"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/madflojo/tasks"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	var (
		db  database.Database
		b   *rod.Browser
		l   *launcher.Launcher
		err error
	)

	defer func() {
		if r := recover(); r != nil {
			slog.Warn(fmt.Sprintf("recovered from panic: %s", r))
		}
		if db != nil {
			if err = db.Close(); err != nil {
				slog.Warn(fmt.Sprintf("error during db conn close: %s", err))
			}
		}
		if b != nil {
			if err = b.Close(); err != nil {
				slog.Warn(fmt.Sprintf("error during browser close: %s", err))
			}
		}
		if l != nil {
			if l != nil {
				l.Cleanup()
			}
		}
		slog.Info("bye")
	}()

	// init DB
	db, err = database.NewDatabase()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("established database connection")

	// init browser
	b, l, err = browser.New()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("initialized browser")

	// prepare
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	slog.Info("starting update...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := tasks.New()
	if err = scheduler.AddWithID("update", task.Update(ctx, db, b)); err != nil {
		slog.Error(fmt.Sprintf("failed to add task: %s", err.Error()))
		return
	}

	<-exit
	cancel()
	scheduler.Stop()
	slog.Info("stopped scheduler. exiting...")
}
