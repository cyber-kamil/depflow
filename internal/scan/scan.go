package scan

import (
	"os"
	"path/filepath"
)

var SupportedLockFiles = []string{
	"go.mod",
	"requirements.txt",
	"Pipfile.lock",
	"pom.xml",
	"build.gradle",
	"package-lock.json",
	"yarn.lock",
}

func ScanForLockFiles(dir string) ([]string, error) {
	found := []string{}
	for _, lf := range SupportedLockFiles {
		path := filepath.Join(dir, lf)
		if _, err := os.Stat(path); err == nil {
			found = append(found, lf)
		}
	}
	return found, nil
}
