# Development Guide

This guide explains how to set up a local development environment for `zcp-ddns`.

## Prerequisites

| Tool   | Minimum Version | Notes                      |
| ------ | --------------- | -------------------------- |
| Go     | 1.25            | Declared in `go.mod`       |
| Make   | Any             | Used for common tasks      |
| Docker | Any current     | Used for container builds  |
| Git    | Any             | Used for version embedding |

## Repository Structure

```text
zcp-ddns/
├── cmd/
│   └── zcp-ddns/        # Service entry point
├── internal/
│   └── version/         # Build-time version string
├── docs/                # Project documentation
├── .github/             # CI, issue templates, CODEOWNERS
├── Dockerfile           # Container image definition
├── Makefile             # Build, test, and quality targets
└── go.mod               # Module definition
```

## Local Workflow

```bash
make fmt          # Format Go and Markdown files
make fmt-check    # Check formatting without writing
make vet          # Run go vet
make test         # Run tests
make build        # Build bin/zcp-ddns
```

## Docker Build

```bash
make docker
```

The image is tagged as:

```text
ghcr.io/zsoftly/zcp-ddns:<git-version>
ghcr.io/zsoftly/zcp-ddns:latest
```

## Testing Approach

Prefer tests around behavior:

- Config parsing.
- IP detection fallback behavior.
- DNS record comparison.
- API client integration boundaries.
- Retry and error handling.

Use fake HTTP servers for API tests. Do not call production ZCP APIs from unit tests.

## Pull Requests

Follow `CONTRIBUTING.md`. Keep PRs tied to one issue and keep decisions in GitHub issues so future contributors can follow the reasoning.
