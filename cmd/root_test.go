package cmd

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/cyber-kamil/depflow/internal/model"
	"github.com/cyber-kamil/depflow/internal/report"
)

func TestWriteMarkdownReportWithHeader(t *testing.T) {
	tmp, err := os.CreateTemp("", "report-*.md")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	reports := []report.NpmDepReport{
		{Name: "express", Current: "4.18.2", Latest: "4.18.2", Outdated: false},
		{Name: "lodash", Current: "4.17.20", Latest: "4.17.21", Outdated: true},
	}
	err = writeMarkdownReportWithHeader("NPM", reports, map[string]*model.ChangelogInfo{}, tmp.Name(), false)
	if err != nil {
		t.Fatalf("failed to write report: %v", err)
	}
	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatalf("failed to read report: %v", err)
	}
	if !strings.Contains(string(data), "express") || !strings.Contains(string(data), "lodash") {
		t.Error("report missing dependency names")
	}
	if !strings.Contains(string(data), "NPM") {
		t.Error("report missing header")
	}
}

func TestScanForNpmLockFile(t *testing.T) {
	dir := t.TempDir()
	f, err := os.Create(dir + "/package-lock.json")
	if err != nil {
		t.Fatalf("failed to create lock file: %v", err)
	}
	f.Close()
	files, _ := ioutil.ReadDir(dir)
	if len(files) == 0 {
		t.Fatalf("no files in temp dir: %s", dir)
	}
	path, err := scanForNpmLockFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(path, "package-lock.json") {
		t.Errorf("expected package-lock.json, got %s (dir: %s, files: %v)", path, dir, files)
	}
}

func TestScanForYarnLockFile(t *testing.T) {
	dir := t.TempDir()
	f, err := os.Create(dir + "/yarn.lock")
	if err != nil {
		t.Fatalf("failed to create lock file: %v", err)
	}
	f.Close()
	files, _ := ioutil.ReadDir(dir)
	if len(files) == 0 {
		t.Fatalf("no files in temp dir: %s", dir)
	}
	path, err := scanForYarnLockFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(path, "yarn.lock") {
		t.Errorf("expected yarn.lock, got %s (dir: %s, files: %v)", path, dir, files)
	}
}

func TestCheckGoDependencies(t *testing.T) {
	dir := t.TempDir()
	gomod := `module github.com/example/project

go 1.20

require (
	github.com/stretchr/testify v1.8.0
)
`
	gomodPath := dir + "/go.mod"
	if err := os.WriteFile(gomodPath, []byte(gomod), 0644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	mockVersionChecker := func(dir string) map[string]string {
		return map[string]string{"github.com/stretchr/testify": "v1.9.0"}
	}

	reports, changelogs, err := checkGoDependencies(dir, gomodPath, func(dir string) (map[string]string, error) {
		return mockVersionChecker(dir), nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reports) == 0 || !reports[0].Outdated {
		t.Errorf("expected outdated module, got %v", reports)
	}
	_ = changelogs // not checked in this test
}
