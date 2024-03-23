package database

import (
	"context"
	"github.com/echovisionlab/aws-weather-updater/internal/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	t.Run("database must return properly", func(t *testing.T) {
		container := testutil.SetupPostgres(context.Background(), t)
		defer testutil.ShutdownContainer(context.Background(), t, container)
		sqlx, err := NewDatabase()
		assert.NotNil(t, sqlx)
		assert.NoError(t, err)
	})
	t.Run("database must return error", func(t *testing.T) {
		testutil.ResetDatabaseEnv(t)
		sqlx, err := NewDatabase()
		assert.Error(t, err)
		assert.Nil(t, sqlx)
		assert.Contains(t, err.Error(), "missing")
	})
}
