package testutil

import (
	"os"
	"path"
)

func GetProjectRoot() string {
	d, _ := os.Getwd()
	return ProjectRootPattern.FindStringSubmatch(d)[1]
}

var TestDataPath = path.Join(GetProjectRoot(), "test", "data")
