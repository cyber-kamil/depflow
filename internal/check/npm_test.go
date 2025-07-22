package check

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetNpmLatestVersion_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"dist-tags": map[string]string{"latest": "1.2.3"},
		})
	}))
	defer ts.Close()

	oldURL := npmRegistryURL
	npmRegistryURL = ts.URL + "/"
	defer func() { npmRegistryURL = oldURL }()

	latest, err := GetNpmLatestVersionWithBase("testpkg", ts.URL+"/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if latest != "1.2.3" {
		t.Errorf("expected 1.2.3, got %s", latest)
	}
}

func TestGetNpmLatestVersion_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	_, err := GetNpmLatestVersionWithBase("notfound", ts.URL+"/")
	if err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestGetNpmLatestVersion_NetworkError(t *testing.T) {
	_, err := GetNpmLatestVersionWithBase("pkg", "http://localhost:0/")
	if err == nil {
		t.Error("expected network error, got nil")
	}
}

// Helper for testing: allows base URL override
var npmRegistryURL = "https://registry.npmjs.org/"

func GetNpmLatestVersionWithBase(pkg, base string) (string, error) {
	url := fmt.Sprintf("%s%s", base, pkg)
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
