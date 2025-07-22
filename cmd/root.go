package cmd

import (
	"fmt"
	"os"

	"github.com/cyber-kamil/depflow/internal/check"
	"github.com/cyber-kamil/depflow/internal/model"
	"github.com/cyber-kamil/depflow/internal/parse"
	"github.com/cyber-kamil/depflow/internal/report"
	"github.com/cyber-kamil/depflow/internal/scan"
	"github.com/spf13/cobra"
)

var (
	dir    string
	output string
)

func scanForNpmLockFile(dir string) (string, error) {
	found, err := scan.ScanForLockFiles(dir)
	if err != nil {
		return "", err
	}
	for _, lf := range found {
		if lf == "package-lock.json" {
			return dir + "/" + lf, nil
		}
	}
	return "", nil
}

func checkNpmDependencies(lockPath string) ([]report.NpmDepReport, map[string]*model.ChangelogInfo, error) {
	deps, err := parse.ParseNpmLockFile(lockPath)
	if err != nil {
		return nil, nil, err
	}
	reports := []report.NpmDepReport{}
	changelogs := make(map[string]*model.ChangelogInfo)
	for name, current := range deps {
		latest, err := check.GetNpmLatestVersion(name)
		if err != nil {
			reports = append(reports, report.NpmDepReport{
				Name:     name,
				Current:  current,
				Latest:   "",
				Outdated: false,
			})
			continue
		}
		outdated := current != latest
		reports = append(reports, report.NpmDepReport{
			Name:     name,
			Current:  current,
			Latest:   latest,
			Outdated: outdated,
		})
		if outdated {
			info, err := check.FetchChangelogInfo(name, current, latest)
			if err == nil && info != nil {
				changelogs[name] = info
			}
		}
	}
	return reports, changelogs, nil
}

func scanForYarnLockFile(dir string) (string, error) {
	found, err := scan.ScanForLockFiles(dir)
	if err != nil {
		return "", err
	}
	for _, lf := range found {
		if lf == "yarn.lock" {
			return dir + "/" + lf, nil
		}
	}
	return "", nil
}

func checkYarnDependencies(lockPath string) ([]report.NpmDepReport, map[string]*model.ChangelogInfo, error) {
	deps, err := parse.ParseYarnLockFile(lockPath)
	if err != nil {
		return nil, nil, err
	}
	reports := []report.NpmDepReport{}
	changelogs := make(map[string]*model.ChangelogInfo)
	for name, current := range deps {
		latest, err := check.GetNpmLatestVersion(name)
		if err != nil {
			reports = append(reports, report.NpmDepReport{
				Name:     name,
				Current:  current,
				Latest:   "",
				Outdated: false,
			})
			continue
		}
		outdated := current != latest
		reports = append(reports, report.NpmDepReport{
			Name:     name,
			Current:  current,
			Latest:   latest,
			Outdated: outdated,
		})
		if outdated {
			info, err := check.FetchChangelogInfo(name, current, latest)
			if err == nil && info != nil {
				changelogs[name] = info
			}
		}
	}
	return reports, changelogs, nil
}

// Add a type for the Go version checker function
// This allows us to inject a mock in tests

type GoVersionChecker func(dir string) (map[string]string, error)

func checkGoDependencies(dir string, goModPath string, versionChecker GoVersionChecker) ([]report.NpmDepReport, map[string]*model.ChangelogInfo, error) {
	mods, err := parse.ParseGoModFile(goModPath)
	if err != nil {
		return nil, nil, err
	}
	latest, err := versionChecker(dir)
	if err != nil {
		return nil, nil, err
	}
	reports := []report.NpmDepReport{}
	changelogs := make(map[string]*model.ChangelogInfo)
	for name, current := range mods {
		newest, hasUpdate := latest[name]
		outdated := hasUpdate && current != newest
		reports = append(reports, report.NpmDepReport{
			Name:     name,
			Current:  current,
			Latest:   newest,
			Outdated: outdated,
		})
		if outdated {
			info, err := check.FetchChangelogInfo(name, current, newest)
			if err == nil && info != nil {
				changelogs[name] = info
			}
		}
	}
	return reports, changelogs, nil
}

func writeMarkdownReport(reports []report.NpmDepReport, output string) error {
	md := report.GenerateNpmMarkdownReport(reports, map[string]*model.ChangelogInfo{})
	return os.WriteFile(output, []byte(md), 0644)
}

func writeMarkdownReportWithHeader(header string, reports []report.NpmDepReport, changelogs map[string]*model.ChangelogInfo, output string, appendMode bool) error {
	md := "## " + header + "\n" + report.GenerateNpmMarkdownReport(reports, changelogs) + "\n"
	flag := os.O_CREATE | os.O_WRONLY
	if appendMode {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}
	f, err := os.OpenFile(output, flag, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(md)
	return err
}

var rootCmd = &cobra.Command{
	Use:   "depflow",
	Short: "Check for outdated dependencies in Go, Python, and Java projects",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("depflow: Scanning directory %s, will output to %s\n", dir, output)

		lockPath, err := scanForNpmLockFile(dir)
		if err != nil {
			fmt.Printf("Error scanning for lock files: %v\n", err)
			return
		}
		yarnPath, err := scanForYarnLockFile(dir)
		if err != nil {
			fmt.Printf("Error scanning for lock files: %v\n", err)
			return
		}

		wrote := false
		if lockPath != "" {
			fmt.Printf("Found lock file: %s\n", lockPath)
			reports, changelogs, err := checkNpmDependencies(lockPath)
			if err != nil {
				fmt.Printf("Error checking npm dependencies: %v\n", err)
				return
			}
			fmt.Println("Checking npm dependencies for updates...")
			err = writeMarkdownReportWithHeader("NPM (package-lock.json)", reports, changelogs, output, false)
			if err != nil {
				fmt.Printf("Error writing report to %s: %v\n", output, err)
			} else {
				fmt.Printf("Report written to %s\n", output)
				wrote = true
			}
		}
		if yarnPath != "" {
			fmt.Printf("Found lock file: %s\n", yarnPath)
			reports, changelogs, err := checkYarnDependencies(yarnPath)
			if err != nil {
				fmt.Printf("Error checking yarn dependencies: %v\n", err)
				return
			}
			fmt.Println("Checking yarn dependencies for updates...")
			err = writeMarkdownReportWithHeader("Yarn (yarn.lock)", reports, changelogs, output, wrote)
			if err != nil {
				fmt.Printf("Error writing report to %s: %v\n", output, err)
			} else {
				fmt.Printf("Report written to %s\n", output)
			}
		}
		goModPath := dir + "/go.mod"
		if _, err := os.Stat(goModPath); err == nil {
			fmt.Printf("Found lock file: %s\n", goModPath)
			reports, changelogs, err := checkGoDependencies(dir, goModPath, check.GetGoModuleLatestVersions)
			if err != nil {
				fmt.Printf("Error checking Go dependencies: %v\n", err)
				return
			}
			fmt.Println("Checking Go dependencies for updates...")
			err = writeMarkdownReportWithHeader("Go (go.mod)", reports, changelogs, output, false)
			if err != nil {
				fmt.Printf("Error writing report to %s: %v\n", output, err)
			} else {
				fmt.Printf("Report written to %s\n", output)
				wrote = true
			}
		}
		if lockPath == "" && yarnPath == "" {
			fmt.Println("No package-lock.json or yarn.lock found.")
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVar(&dir, "dir", ".", "Directory to scan for dependency files")
	rootCmd.PersistentFlags().StringVar(&output, "output", "dependency-report.md", "Output Markdown report file")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
