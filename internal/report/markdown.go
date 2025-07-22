package report

import "github.com/cyber-kamil/depflow/internal/model"

type NpmDepReport struct {
	Name     string
	Current  string
	Latest   string
	Outdated bool
}

// GenerateNpmMarkdownReport generates a Markdown report for npm dependencies, including changelog links and highlights if provided.
func GenerateNpmMarkdownReport(deps []NpmDepReport, changelogs map[string]*model.ChangelogInfo) string {
	report := "# NPM Dependency Update Report\n\n"
	report += "| Dependency | Current Version | Latest Version | Status | Changelog | Highlights |\n"
	report += "|------------|-----------------|---------------|--------|-----------|------------|\n"
	for _, dep := range deps {
		status := "Up to date"
		if dep.Outdated {
			status = "Update available"
		}
		changelog := ""
		highlights := ""
		if info, ok := changelogs[dep.Name]; ok {
			if info.ChangelogURL != "" {
				changelog = "[Changelog](" + info.ChangelogURL + ")"
			}
			if len(info.Highlights) > 0 {
				highlights = "- " + info.Highlights[0]
				for _, h := range info.Highlights[1:] {
					highlights += "<br>- " + h
				}
			}
		}
		report += "| " + dep.Name + " | " + dep.Current + " | " + dep.Latest + " | " + status + " | " + changelog + " | " + highlights + " |\n"
	}
	return report
}
