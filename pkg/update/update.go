package update

import (
	"context"
	"fmt"
	"github.com/echovisionlab/aws-weather-updater/pkg/database"
	"github.com/echovisionlab/aws-weather-updater/pkg/model"
)

func bindThenRun(ctx context.Context, db database.Database, query string, arg interface{}) (int64, error) {
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

func ClearBefore(ctx context.Context, db database.Database, days int) (int64, error) {
	q := fmt.Sprintf(DeleteOlderThanThreeDaysQuery, days)
	if n, err := db.ExecContext(ctx, q); err != nil {
		return 0, fmt.Errorf("failed to delete rows older than %d days: %w", days, err)
	} else if rowsAffected, err := n.RowsAffected(); err != nil {
		return 0, fmt.Errorf("failed to delete rows older than %d days: %w", days, err)
	} else {
		return rowsAffected, nil
	}
}

func Records(ctx context.Context, db database.Database, records []model.Record) (int64, error) {
	if n, err := bindThenRun(ctx, db, UpsertRecordQuery, records); err != nil {
		return 0, fmt.Errorf("failed to update records: %w", err)
	} else {
		return n, nil
	}
}

func Stations(ctx context.Context, db database.Database, stations []model.Station) (int64, error) {
	if n, err := bindThenRun(ctx, db, UpsertStationQuery, stations); err != nil {
		return 0, fmt.Errorf("failed to update stations: %w", err)
	} else {
		return n, nil
	}
}
