package parse

import (
	"os"
	"testing"
)

func TestParseNpmLockFile_Valid(t *testing.T) {
	json := `{"dependencies": {"express": {"version": "4.18.2"}, "lodash": {"version": "4.17.21"}}}`
	tmp, err := os.CreateTemp("", "package-lock-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write([]byte(json)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmp.Close()

	deps, err := ParseNpmLockFile(tmp.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deps["express"] != "4.18.2" || deps["lodash"] != "4.17.21" {
		t.Errorf("unexpected deps: %v", deps)
	}
}

func TestParseNpmLockFile_MissingFile(t *testing.T) {
	_, err := ParseNpmLockFile("nonexistent.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParseNpmLockFile_Malformed(t *testing.T) {
	json := `{"dependencies": {"express": {"version": 4.18.2}}}` // version not a string
	tmp, err := os.CreateTemp("", "package-lock-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write([]byte(json)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmp.Close()

	_, err = ParseNpmLockFile(tmp.Name())
	if err == nil {
		t.Error("expected error for malformed file, got nil")
	}
}
