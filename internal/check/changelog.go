package check

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/cyber-kamil/depflow/internal/model"
)

// FetchChangelogInfo tries to find and summarize the changelog for a dependency.
// For now, only npm is supported. This function can be extended for other ecosystems.
func FetchChangelogInfo(depName, currentVersion, latestVersion string) (*model.ChangelogInfo, error) {
	// Step 1: Fetch npm package metadata
	url := fmt.Sprintf("https://registry.npmjs.org/%s", depName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch npm info for %s: %w", depName, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("npm registry returned status %d for %s", resp.StatusCode, depName)
	}
	var data struct {
		Repository struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"repository"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode npm registry response for %s: %w", depName, err)
	}

	repoURL := data.Repository.URL
	if repoURL == "" {
		return &model.ChangelogInfo{
			Dependency:   depName,
			RepoURL:      "",
			ChangelogURL: "",
			Highlights:   nil,
		}, nil
	}

	// Normalize GitHub URLs (remove git+, .git, etc.)
	repoURL = strings.TrimPrefix(repoURL, "git+")
	repoURL = strings.TrimSuffix(repoURL, ".git")
	repoURL = strings.Replace(repoURL, "://github.com/", "://github.com/", 1)
	if strings.HasPrefix(repoURL, "git@github.com:") {
		repoURL = "https://github.com/" + strings.TrimPrefix(repoURL, "git@github.com:")
	}

	// Try to construct a likely changelog URL
	changelogURL := repoURL + "/blob/master/CHANGELOG.md"
	// Try to fetch and summarize changelog content if GitHub
	highlights := []string{}
	if strings.Contains(repoURL, "github.com/") {
		ownerRepo := strings.TrimPrefix(repoURL, "https://github.com/")
		ownerRepo = strings.TrimSuffix(ownerRepo, "/")
		branches := []string{"main", "master"}
		var changelogContent string
		for _, branch := range branches {
			rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/CHANGELOG.md", ownerRepo, branch)
			resp, err := http.Get(rawURL)
			if err == nil && resp.StatusCode == 200 {
				b, _ := io.ReadAll(resp.Body)
				changelogContent = string(b)
				changelogURL = rawURL
				resp.Body.Close()
				break
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
		if changelogContent != "" {
			// Extract section between currentVersion and latestVersion
			section := extractChangelogSection(changelogContent, currentVersion, latestVersion)
			highlights = extractBreakingHighlights(section)
		}
	}

	return &model.ChangelogInfo{
		Dependency:   depName,
		RepoURL:      repoURL,
		ChangelogURL: changelogURL,
		Highlights:   highlights,
	}, nil
}

// extractChangelogSection extracts the changelog section(s) between current and latest version headings.
func extractChangelogSection(content, current, latest string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var lines []string
	inSection := false
	verRe := regexp.MustCompile(`^##? ?v?([0-9]+\.[0-9]+\.[0-9]+)`) // e.g. ## 4.17.21 or # v4.17.21
	foundLatest := false
	for scanner.Scan() {
		line := scanner.Text()
		if verRe.MatchString(line) {
			ver := verRe.FindStringSubmatch(line)[1]
			if ver == latest {
				inSection = true
				foundLatest = true
				continue
			}
			if ver == current && foundLatest {
				break
			}
		}
		if inSection {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

// extractBreakingHighlights scans changelog text for lines with breaking changes or important keywords.
func extractBreakingHighlights(section string) []string {
	highlights := []string{}
	scanner := bufio.NewScanner(strings.NewReader(section))
	keywords := []string{"break", "major", "deprecat", "remove", "migration"}
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		for _, kw := range keywords {
			if strings.Contains(line, kw) {
				highlights = append(highlights, scanner.Text())
				break
			}
		}
	}
	return highlights
}
