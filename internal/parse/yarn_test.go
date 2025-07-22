package parse

import (
	"os"
	"testing"
)

func TestParseYarnLockFile_Valid(t *testing.T) {
	lock := `lodash@^4.17.20:
  version "4.17.21"
  resolved "https://registry.yarnpkg.com/lodash/-/lodash-4.17.21.tgz"
  integrity sha512-v2kDE...

express@^4.18.2:
  version "4.18.2"
  resolved "https://registry.yarnpkg.com/express/-/express-4.18.2.tgz"
  integrity sha512-abc...`
	tmp, err := os.CreateTemp("", "yarn-*.lock")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write([]byte(lock)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmp.Close()

	deps, err := ParseYarnLockFile(tmp.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deps["lodash"] != "4.17.21" || deps["express"] != "4.18.2" {
		t.Errorf("unexpected deps: %v", deps)
	}
}

func TestParseYarnLockFile_MissingFile(t *testing.T) {
	_, err := ParseYarnLockFile("nonexistent.lock")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParseYarnLockFile_Malformed(t *testing.T) {
	lock := `lodash@^4.17.20:
  version 4.17.21` // version not quoted
	tmp, err := os.CreateTemp("", "yarn-*.lock")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write([]byte(lock)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmp.Close()

	deps, err := ParseYarnLockFile(tmp.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 0 {
		t.Errorf("expected no deps, got %v", deps)
	}
}
