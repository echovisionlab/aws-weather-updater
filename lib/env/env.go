package env

import (
	"fmt"
	"github.com/echovisionlab/aws-weather-updater/lib/constants"
	"os"
	"strings"
)

func CheckRequiredEnv() error {
	missing := make([]string, 0)

	for _, key := range []string{constants.DatabaseName, constants.DatabaseUser, constants.DatabasePass} {
		if v := os.Getenv(key); len(v) == 0 {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: [%+v]", strings.Join(missing, ", "))
	}

	return nil
}
