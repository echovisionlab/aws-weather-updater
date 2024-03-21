package update

import (
	"context"
	"fmt"
	"github.com/echovisionlab/aws-weather-updater/lib/database"
	"github.com/echovisionlab/aws-weather-updater/lib/model"
	"github.com/echovisionlab/aws-weather-updater/lib/query"
	"github.com/echovisionlab/aws-weather-updater/lib/types/fetchresult"
	"sync"
)

func doUpdate(ctx context.Context, db database.Database, query string, arg interface{}) (int64, error) {
	q, args, err := db.BindNamed(query, arg)
	if err != nil {
		return 0, fmt.Errorf("failed to bind named query: %w", err)
	}

	result, err := db.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute with context: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to extract rows affected: %w", err)
	}

	return rows, err
}

// Run updates both stations and records.
// It will fail fast if given context is either cancelled, timed out.
// It will also fail fast if anything goes wrong during the process it executes.
func Run(ctx context.Context, wg *sync.WaitGroup, db database.Database, items fetchresult.FetchResult) error {
	defer wg.Done()
	if _, err := updateStations(ctx, db, items.Stations()); err != nil {
		return fmt.Errorf("failed to update stations: %w", err)
	}
	if _, err := updateRecords(ctx, db, items.Records()); err != nil {
		return fmt.Errorf("failed to update records: %w", err)
	}
	return nil
}

func updateRecords(ctx context.Context, db database.Database, records []model.Record) (int64, error) {
	return doUpdate(ctx, db, query.UpsertRecordQuery, records)
}

func updateStations(ctx context.Context, db database.Database, stations []model.Station) (int64, error) {
	return doUpdate(ctx, db, query.UpsertStationQuery, stations)
}
