package parse

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/mod/modfile"
)

// ParseGoModFile parses a go.mod file and returns a map of module names to their versions.
func ParseGoModFile(path string) (map[string]string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read go.mod: %w", err)
	}
	mf, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse go.mod: %w", err)
	}
	mods := make(map[string]string)
	for _, req := range mf.Require {
		mods[req.Mod.Path] = req.Mod.Version
	}
	return mods, nil
}
