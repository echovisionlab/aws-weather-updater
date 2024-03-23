package env

import (
	"fmt"
	"os"
	"strings"
)

func CheckRequiredEnv() error {
	missing := make([]string, 0)

	for _, key := range []string{DatabaseName, DatabaseUser, DatabasePass} {
		if v := os.Getenv(key); len(v) == 0 {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: [%+v]", strings.Join(missing, ", "))
	}

	return nil
}
