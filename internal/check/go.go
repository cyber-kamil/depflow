package check

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type GoModuleVersion struct {
	Path    string
	Version string
	Update  *struct {
		Version string
	} `json:"Update"`
}

// GetGoModuleLatestVersions runs 'go list -m -u -json all' in the given directory and returns a map of module names to their latest versions (if available).
func GetGoModuleLatestVersions(dir string) (map[string]string, error) {
	cmd := exec.Command("go", "list", "-m", "-u", "-json", "all")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run 'go list': %w", err)
	}

	latest := make(map[string]string)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(splitJSONObjects)
	for scanner.Scan() {
		var mod GoModuleVersion
		if err := json.Unmarshal(scanner.Bytes(), &mod); err != nil {
			continue // skip malformed
		}
		if mod.Update != nil && mod.Update.Version != "" {
			latest[mod.Path] = mod.Update.Version
		}
	}
	return latest, nil
}

// splitJSONObjects is a bufio.SplitFunc that splits a stream of concatenated JSON objects.
func splitJSONObjects(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var depth, start int
	for i, b := range data {
		switch b {
		case '{':
			if depth == 0 {
				start = i
			}
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i + 1, data[start : i+1], nil
			}
		}
	}
	if atEOF && depth == 0 && start < len(data) {
		return len(data), data[start:], nil
	}
	return 0, nil, nil
}
