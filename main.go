package main

import (
	"context"
	"github.com/echovisionlab/aws-weather-updater/pkg/database"
	"github.com/echovisionlab/aws-weather-updater/pkg/fetch"
	"github.com/echovisionlab/aws-weather-updater/pkg/type/fetchresult"
	"github.com/echovisionlab/aws-weather-updater/pkg/update"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// init DB
	db, err := database.NewDatabase()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer closeDatabase(db)
	slog.Info("established database connection")

	// prepare
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	fetchResultChan := make(chan fetchresult.FetchResult)
	errorChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())

	slog.Info("starting update...")
	// run
	wg.Add(1)
	go func() {
		defer wg.Done()
		if r, err := fetch.StationsAndRecords(ctx, time.Now()); err != nil {
			errorChan <- err
		} else {
			slog.Info("fetched new data. updating...")
			fetchResultChan <- r
		}
	}()

	select {
	case <-exit:
		slog.Info("operation cancelled")
		cancel()
	case r := <-fetchResultChan:
		wg.Add(1)
		if err := update.Run(ctx, &wg, db, r); err != nil { // blocking intentionally
			slog.Error(err.Error())
		} else {
			slog.Info("updated stations and records")
		}
	case err := <-errorChan:
		slog.Error(err.Error())
		cancel()
	}

	wg.Wait()
	slog.Info("goodbye...")
}

func closeDatabase(db database.Database) {
	if err := db.Close(); err != nil {
		slog.Error("error during database connection close: %w", err)
	}
}
