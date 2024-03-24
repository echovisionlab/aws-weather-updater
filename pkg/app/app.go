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
	"strconv"
	"time"
)

func Run(exit <-chan os.Signal) {
	var (
		db  database.Database
		b   *rod.Browser
		l   *launcher.Launcher
		err error
	)

	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()
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

	interval, err := getInterval()
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info("starting update...")

	scheduler := tasks.New()
	if err = scheduler.AddWithID("update", task.Update(ctx, db, b, interval)); err != nil {
		slog.Error(fmt.Sprintf("failed to add task: %s", err.Error()))
		return
	}

	<-exit
}

func getInterval() (time.Duration, error) {
	s := os.Getenv("INTERVAL_SECONDS")
	if len(s) == 0 {
		slog.Info("missing INTERVAL_SECONDS environment. setting to default: 1 min")
		return time.Minute, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return *new(time.Duration), fmt.Errorf("invalid interval_seconds format: %s", s)
	}
	return time.Second * time.Duration(v), nil
}
