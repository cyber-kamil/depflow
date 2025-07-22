package parse

import (
	"os"
	"testing"
)

func TestParseGoModFile_Valid(t *testing.T) {
	gomod := `module github.com/example/project

go 1.20

require (
	github.com/stretchr/testify v1.8.0
	golang.org/x/mod v0.12.0
)
`
	tmp, err := os.CreateTemp("", "go-*.mod")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write([]byte(gomod)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmp.Close()

	mods, err := ParseGoModFile(tmp.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mods["github.com/stretchr/testify"] != "v1.8.0" || mods["golang.org/x/mod"] != "v0.12.0" {
		t.Errorf("unexpected mods: %v", mods)
	}
}

func TestParseGoModFile_MissingFile(t *testing.T) {
	_, err := ParseGoModFile("nonexistent.mod")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParseGoModFile_Malformed(t *testing.T) {
	gomod := `module github.com/example/project
require github.com/stretchr/testify` // missing version
	tmp, err := os.CreateTemp("", "go-*.mod")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write([]byte(gomod)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmp.Close()

	_, err = ParseGoModFile(tmp.Name())
	if err == nil {
		t.Error("expected error for malformed file, got nil")
	}
}
