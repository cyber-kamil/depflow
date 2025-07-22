package cmd

import (
	"os"
	"strings"
	"testing"

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
	err = writeMarkdownReportWithHeader("NPM", reports, tmp.Name(), false)
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
	path, err := scanForNpmLockFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(path, "package-lock.json") {
		t.Errorf("expected package-lock.json, got %s", path)
	}
}

func TestScanForYarnLockFile(t *testing.T) {
	dir := t.TempDir()
	f, err := os.Create(dir + "/yarn.lock")
	if err != nil {
		t.Fatalf("failed to create lock file: %v", err)
	}
	f.Close()
	path, err := scanForYarnLockFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(path, "yarn.lock") {
		t.Errorf("expected yarn.lock, got %s", path)
	}
}
