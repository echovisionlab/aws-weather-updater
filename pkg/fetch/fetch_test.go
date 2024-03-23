package fetch

import (
	"context"
	"github.com/echovisionlab/aws-weather-updater/pkg/browser"
	"github.com/go-rod/rod/lib/proto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_StationsAndRecords(t *testing.T) {
	b, l, err := browser.New()

	cleanup := func() {
		_ = b.Close()
		l.Cleanup()
		t.FailNow()
	}
	if err != nil {
		cleanup()
	}

	p, err := b.Page(proto.TargetCreateTarget{
		URL:                     "",
		Width:                   nil,
		Height:                  nil,
		BrowserContextID:        "",
		EnableBeginFrameControl: false,
		NewWindow:               false,
		Background:              false,
		ForTab:                  false,
	})
	if err != nil {
		_ = p.Close()
		cleanup()
	}

	t.Run("must report context exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer cancel()
		result, err := StationsAndRecords(p, ctx, time.Now())
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
	t.Run("must report context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		result, err := StationsAndRecords(p, ctx, time.Now())
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, context.Canceled)
	})
	t.Run("must fetch", func(t *testing.T) {
		result, err := StationsAndRecords(p, context.Background(), time.Now())
		assert.NoError(t, err)
		assert.NotEmpty(t, result.Records())
		assert.NotEmpty(t, result.Stations())
	})
}
