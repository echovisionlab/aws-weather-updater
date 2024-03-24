package fetch

import (
	"context"
	"github.com/echovisionlab/aws-weather-updater/pkg/browser"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_StationsAndRecords(t *testing.T) {
	b, l, err := browser.New()

	if err != nil {
		t.Fail()
		return
	}

	t.Cleanup(func() {
		if b != nil {
			b.MustClose()
		}
		if l != nil {
			l.Cleanup()
		}
	})

	t.Run("must report context exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer cancel()

		p := b.MustPage()
		defer p.MustClose()
		result, err := StationsAndRecords(ctx, p, time.Now())
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
	t.Run("must report context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		p := b.MustPage()
		defer p.MustClose()
		result, err := StationsAndRecords(ctx, p, time.Now())
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, context.Canceled)
	})
	t.Run("must fetch", func(t *testing.T) {
		p := b.MustPage()
		defer p.MustClose()
		result, err := StationsAndRecords(context.Background(), p, time.Now())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		if t.Failed() {
			return
		}
		assert.NotNil(t, result.Stations())
		assert.NotNil(t, result.Records())
		if t.Failed() {
			return
		}
		assert.NotEmpty(t, result.Stations())
		assert.NotEmpty(t, result.Records())
	})
}
