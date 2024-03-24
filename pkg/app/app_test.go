package app

import (
	"context"
	"github.com/echovisionlab/aws-weather-updater/internal/testutil"
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	// prep
	postgres := testutil.SetupPostgres(context.Background(), t)
	defer testutil.ShutdownContainer(context.Background(), t, postgres)

	t.Run("must not panic", func(t *testing.T) {
		t.Run("operation cancelled", func(t *testing.T) {
			exit := make(chan os.Signal)
			go func() {
				<-time.After(time.Millisecond * 100)
				exit <- syscall.SIGINT
			}()
			assert.NotPanics(t, func() {
				Run(exit)
			})
		})
	})
}

func Test_getInterval(t *testing.T) {
	key := "INTERVAL_SECONDS"

	defer t.Setenv(key, "")
	set := func(v string) { t.Setenv(key, v) }

	t.Run("must return error", func(t *testing.T) {
		t.Run("when string", func(t *testing.T) {
			set("invalid")
			v, err := getInterval()
			assert.Error(t, err)
			assert.ErrorContains(t, err, "invalid")
			assert.Equal(t, v, *new(time.Duration))
		})
	})

	t.Run("must return value", func(t *testing.T) {
		t.Run("when empty", func(t *testing.T) {
			set("")
			v, err := getInterval()
			assert.NoError(t, err)
			assert.Equal(t, v, time.Minute)
		})
		t.Run("when number", func(t *testing.T) {
			set("10")
			v, err := getInterval()
			assert.NoError(t, err)
			assert.Equal(t, v, time.Second*10)
		})
	})
}
