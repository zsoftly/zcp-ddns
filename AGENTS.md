# Repository Instructions

This repository is the public home for `zcp-ddns`, the ZSoftly Cloud Platform Dynamic DNS service.

## Project Scope

`zcp-ddns` is a small Go service that runs as a container, detects public IP changes, and creates or updates DNS records through ZSoftly DNS.

Keep the project focused:

- One service binary.
- Docker-first runtime.
- DNS-scoped ZCP API access only.
- No broad platform token requirements.
- No internal ZSoftly API shortcuts.

## Source of Truth

Use these files before making changes:

- `README.md` for project purpose and user-facing status.
- `docs/architecture.md` for design boundaries.
- `docs/configuration.md` for planned configuration shape.
- `docs/development.md` for local workflow.
- `docs/roadmap.md` for planned work.
- `CONTRIBUTING.md` for PR rules.

## Code Rules

- Use Go standard library first.
- Keep dependencies minimal.
- Route all ZCP DNS API access through the official ZCP client from `github.com/zsoftly/zcp-cli/pkg/api/dns`.
- If the official client lacks a method, open or update an issue in `zcp-cli` rather than adding ad hoc HTTP calls here.
- Do not commit secrets, tokens, real customer domains, real customer names, or production DNS zones.
- Use placeholder domains such as `example.com` or assigned contributor test zones only.
- Prefer small PRs tied to one GitHub issue.
- Add or update tests for changed behavior.

## Local Checks

Run these before opening a PR:

```bash
make fmt-check
make vet
make test
make build
```

Run `make fmt` before committing if formatting fails.

## Documentation Rules

- Keep docs concise and practical.
- Do not claim features as shipped until code exists and tests pass.
- Mark planned behavior clearly as planned.
- Use ASCII punctuation in documentation and code.

## Release Rules

- Do not create releases from contributor branches.
- Releases are cut by ZSoftly maintainers only.
- Keep `.github/CODEOWNERS` and `MAINTAINERS.md` aligned.
