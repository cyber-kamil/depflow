package model

type ChangelogInfo struct {
	Dependency   string
	RepoURL      string
	ChangelogURL string
	Highlights   []string // e.g., breaking changes or summary lines
}
