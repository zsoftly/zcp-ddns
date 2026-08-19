# Contributor Onboarding

Welcome. This guide takes you from zero to a merged PR.

## 1. Get a ZCP account

1. Register at <https://cloud.zcp.zsoftly.ca/register> and verify your email.
2. Join the community Slack at <https://zcp.zsoftly.ca/community> and ask to be added to `#zcp-open-source`.
3. Email [community@zsoftly.ca](mailto:community@zsoftly.ca) from your signup address. A maintainer will apply contributor credits to your account and, if you need one, assign you a DNS test zone. Do not post email addresses, tokens, or credentials in public channels.
4. Generate an API token in the ZCP console. Contributor tokens are scoped to DNS operations on your test account.

Never use a production domain or production credentials for development or testing.

## 2. Understand the stack

- ZCP DNS docs: <https://docs.zcp.zsoftly.ca>
- API base: <https://api.zcp.zsoftly.ca/api>
- Go API client: this project uses [`github.com/zsoftly/zcp-cli/pkg/api/dns`](https://github.com/zsoftly/zcp-cli) as its client library. Read that package before writing API code. If the client is missing something you need, open an issue in zcp-cli rather than working around it here.
- For prior art on DNS operations, see the DNS resources in [terraform-provider-zcp](https://github.com/zsoftly/terraform-provider-zcp).

## 3. Set up locally

1. Install Go 1.25+ and Docker.
2. Fork this repo and clone your fork.
3. The repository currently contains scaffolding. Run `make test` to verify the current state. Run `make build` after `cmd/zcp-ddns` is added.

## 4. Pick up work

1. Find an open [issue](https://github.com/zsoftly/zcp-ddns/issues). Issues labeled `good first issue` are scoped for newcomers.
2. Comment on the issue to claim it. One issue at a time until your first PR is merged.
3. If anything about the issue is unclear, ask on the issue itself so the answer is preserved for the next person.

## 5. Submit

1. Branch from `main` in your fork.
2. Keep PRs small and focused on one issue.
3. Sign off every commit: `git commit -s` (DCO required, maintainers verify it before merge).
4. Make sure lint, tests, and build pass locally.
5. Open the PR against `main` and reference the issue it closes.

See [CONTRIBUTING.md](CONTRIBUTING.md) for the review target. After a few merged PRs you may be offered triage access. Each project has at most 3 developer maintainers, listed in [MAINTAINERS.md](MAINTAINERS.md). Sustained contributors are invited as slots open.

## 6. Ground rules

1. License is [Apache-2.0](LICENSE). By contributing you agree your work is licensed the same way.
2. No secrets, tokens, real customer names, or registrable production domains anywhere in code, tests, fixtures, or docs. Use your assigned test zone and placeholder data.
3. Decisions live in GitHub issues. Slack is for discussion.
4. Be respectful. [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) applies everywhere the project operates.
