package fetch

import (
	"context"
	"github.com/echovisionlab/aws-weather-updater/internal/testutil"
	"github.com/echovisionlab/aws-weather-updater/pkg/browser"
	"github.com/go-rod/rod/lib/proto"
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
	"time"
)

func Test_StationsAndRecords(t *testing.T) {
	b, l, err := browser.New()

	if err != nil {
		t.FailNow()
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

	backupUrl := measurementBaseUrl
	measurementBaseUrl = path.Join(testutil.TestDataPath, "test_fetch.html")

	t.Cleanup(func() {
		if p != nil {
			p.MustClose()
		}
		if b != nil {
			b.MustClose()
		}
		if l != nil {
			l.Cleanup()
		}
		measurementBaseUrl = backupUrl
	})

	t.Run("must report context exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer cancel()
		result, err := StationsAndRecords(ctx, p, time.Now())
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
	t.Run("must report context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		result, err := StationsAndRecords(ctx, p, time.Now())
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, context.Canceled)
	})
	t.Run("must fetch", func(t *testing.T) {
		result, err := StationsAndRecords(context.Background(), p, time.Now())
		assert.NoError(t, err)
		assert.NotEmpty(t, result.Records())
		assert.NotEmpty(t, result.Stations())
	})
}
