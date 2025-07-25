package report

import (
	"strings"
	"testing"

	"github.com/cyber-kamil/depflow/internal/model"
)

func TestGenerateNpmMarkdownReport(t *testing.T) {
	deps := []NpmDepReport{
		{Name: "express", Current: "4.18.2", Latest: "4.18.2", Outdated: false},
		{Name: "lodash", Current: "4.17.20", Latest: "4.17.21", Outdated: true},
	}
	report := GenerateNpmMarkdownReport(deps, map[string]*model.ChangelogInfo{})
	if !strings.Contains(report, "express") || !strings.Contains(report, "lodash") {
		t.Error("report missing dependency names")
	}
	if !strings.Contains(report, "4.17.20") && !strings.Contains(report, "Update available") {
		t.Error("report missing update info")
	}
	if !strings.Contains(report, "Up to date") {
		t.Error("report missing up to date info")
	}
}
