# depflow

[![Go Report Card](https://goreportcard.com/badge/github.com/cyber-kamil/depflow)](https://goreportcard.com/report/github.com/cyber-kamil/depflow)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**depflow** is an open-source CLI tool that helps you and your team keep dependencies up-to-date and safe across multiple languages. It scans your project for lock files, checks for newer versions, and generates a beautiful Markdown report with changelogs and breaking change highlights—perfect for local use and CI/CD pipelines.

---

## 🚀 Features
- **Multi-language support:**
  - JavaScript/TypeScript: `package-lock.json` (npm), `yarn.lock` (Yarn)
  - Go: `go.mod`
  - *(Planned: Python, Java, and more!)*
- **Detects outdated dependencies** and shows current/latest versions
- **Fetches and links to changelogs** (GitHub, etc.)
- **Highlights breaking changes** from changelogs
- **Markdown report** for easy review or CI artifacts
- **Ready for CI/CD integration**
- **Extensible**: Easy to add new languages and features

---

## 📦 Installation

### Install with Go (recommended)

```sh
go install github.com/cyber-kamil/depflow@latest
```

Or specify a particular version:

```sh
go install github.com/cyber-kamil/depflow@v0.1.0
```

This will place `depflow` in your `$GOPATH/bin` or `$HOME/go/bin`.

### Clone the repo and build

```sh
git clone https://github.com/cyber-kamil/depflow.git
cd depflow
go build -o depflow
```

Or run directly with Go:

```sh
go run . --dir /path/to/your/project
```

---

## 🛠 Usage

Scan a project and generate a Markdown report:

```sh
./depflow --dir /path/to/your/project --output dependency-report.md
```

- `--dir`   : Directory to scan (default: current directory)
- `--output`: Output Markdown file (default: `dependency-report.md`)

### Example Output

```
# NPM Dependency Update Report

| Dependency | Current Version | Latest Version | Status           | Changelog                | Highlights                |
|------------|-----------------|---------------|------------------|--------------------------|---------------------------|
| lodash     | 4.17.20         | 4.17.21       | Update available | [Changelog](...)         | - breaking: removed ...   |
| express    | 4.18.2          | 4.18.2        | Up to date       |                          |                           |
```

---

## 🤖 CI/CD Integration

Add to your pipeline to fail on breaking changes or just to generate a report:

```yaml
- name: Check dependencies
  run: ./depflow --dir . --output dependency-report.md
```

---

## 🤝 Contributing

We welcome contributions from everyone! Whether you’re fixing bugs, adding features, or improving docs, your help is appreciated.

- **Fork** the repo and create your branch
- **Add tests** for new features
- **Open a pull request** and describe your changes
- **Join the discussion**: Open issues for bugs, feature requests, or questions

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for more details (coming soon).

---

## 🌍 Community & Support
- **Issues:** [Submit here](https://github.com/cyber-kamil/depflow/issues)
- **Discussions:** [Start a topic](https://github.com/cyber-kamil/depflow/discussions) (coming soon)
- **Roadmap:** See below and suggest new features!

---

## 🗺 Roadmap
- Python (`requirements.txt`, `Pipfile.lock`)
- Java (`pom.xml`, `build.gradle`)
- More changelog/release note sources
- JSON/HTML report formats
- More CI/CD integrations
- Community plugins/extensions

---

## 📄 License
MIT — free for personal and commercial use.

---

*Made by the cyber-kamil* 
