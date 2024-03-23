package testutil

import (
	"github.com/echovisionlab/aws-weather-updater/pkg/env"
	"testing"
)

func ResetDatabaseEnv(t *testing.T) {
	t.Setenv(env.DatabaseHost, "")
	t.Setenv(env.DatabasePort, "")
	t.Setenv(env.DatabaseUser, "")
	t.Setenv(env.DatabasePass, "")
	t.Setenv(env.DatabaseName, "")
}
