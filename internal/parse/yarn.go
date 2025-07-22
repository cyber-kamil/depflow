package parse

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ParseYarnLockFile parses a yarn.lock file (Yarn v1/classic) and returns a map of dependency names to their versions.
func ParseYarnLockFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open yarn.lock: %w", err)
	}
	defer file.Close()

	deps := make(map[string]string)
	scanner := bufio.NewScanner(file)
	var currentDep string
	versionRe := regexp.MustCompile(`^  version "([^"]+)"`)
	depRe := regexp.MustCompile(`^([^"]+)@.*:`)

	for scanner.Scan() {
		line := scanner.Text()
		if depRe.MatchString(line) {
			// e.g. "lodash@^4.17.20:", "react@^17.0.0, react@^17.0.1:"
			depLine := strings.Split(line, ":")[0]
			depNames := strings.Split(depLine, ",")
			// Only take the first dep name (main package)
			currentDep = strings.TrimSpace(strings.Split(depNames[0], "@")[0])
		} else if versionRe.MatchString(line) && currentDep != "" {
			matches := versionRe.FindStringSubmatch(line)
			if len(matches) == 2 {
				deps[currentDep] = matches[1]
				currentDep = ""
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan yarn.lock: %w", err)
	}
	return deps, nil
}
