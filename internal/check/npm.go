package check

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetNpmLatestVersion queries the npm registry for the latest version of a package.
func GetNpmLatestVersion(pkg string) (string, error) {
	url := fmt.Sprintf("https://registry.npmjs.org/%s", pkg)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch npm info for %s: %w", pkg, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("npm registry returned status %d for %s", resp.StatusCode, pkg)
	}

	var data struct {
		DistTags struct {
			Latest string `json:"latest"`
		} `json:"dist-tags"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode npm registry response for %s: %w", pkg, err)
	}

	return data.DistTags.Latest, nil
}
