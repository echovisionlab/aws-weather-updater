package fetch

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_StationsAndRecords(t *testing.T) {
	t.Run("must report context exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer cancel()
		result, err := StationsAndRecords(ctx, time.Now())
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
	t.Run("must report context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		result, err := StationsAndRecords(ctx, time.Now())
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, context.Canceled)
	})
	t.Run("must fetch", func(t *testing.T) {
		result, err := StationsAndRecords(context.Background(), time.Now())
		assert.NoError(t, err)
		assert.NotEmpty(t, result.Records())
		assert.NotEmpty(t, result.Stations())
	})
}
