package parse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type NpmDependency struct {
	Version string `json:"version"`
}

type NpmLockFile struct {
	Dependencies map[string]NpmDependency `json:"dependencies"`
}

// ParseNpmLockFile parses a package-lock.json file and returns a map of dependency names to their versions.
func ParseNpmLockFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open package-lock.json: %w", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read package-lock.json: %w", err)
	}

	var lock NpmLockFile
	if err := json.Unmarshal(bytes, &lock); err != nil {
		return nil, fmt.Errorf("failed to parse package-lock.json: %w", err)
	}

	deps := make(map[string]string)
	for name, dep := range lock.Dependencies {
		deps[name] = dep.Version
	}
	return deps, nil
}
