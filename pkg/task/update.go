package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/echovisionlab/aws-weather-updater/pkg/database"
	"github.com/echovisionlab/aws-weather-updater/pkg/fetch"
	"github.com/echovisionlab/aws-weather-updater/pkg/update"
	"github.com/go-rod/rod"
	"github.com/madflojo/tasks"
	"log/slog"
	"time"
)

func Update(ctx context.Context, db database.Database, b *rod.Browser) *tasks.Task {
	return &tasks.Task{
		TaskContext:            tasks.TaskContext{Context: ctx},
		Interval:               time.Minute,
		RunSingleInstance:      true,
		ErrFunc:                nil,
		ErrFuncWithTaskContext: handleErr,
		FuncWithTaskContext:    doUpdate(db, b),
	}
}

func handleErr(taskContext tasks.TaskContext, err error) {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return
	}
	slog.Error(fmt.Sprintf("error during task '%s': %s", taskContext.ID(), err))
}

func doUpdate(db database.Database, b *rod.Browser) func(tasks.TaskContext) error {
	return func(taskContext tasks.TaskContext) error {
		ctx := taskContext.Context
		page := b.MustPage()
		defer page.MustClose()
		fetched, err := fetch.StationsAndRecords(ctx, page, time.Now())
		if err != nil {
			return err
		}
		if _, err = update.Stations(ctx, db, fetched.Stations()); err != nil {
			return err
		}
		if _, err = update.Records(ctx, db, fetched.Records()); err != nil {
			return err
		}
		return nil
	}
}
