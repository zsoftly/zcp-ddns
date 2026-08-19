# Contributing

Thank you for your interest in zcp-ddns. New here? Start with [ONBOARDING.md](ONBOARDING.md) for account setup and test credentials.

## Reporting Issues

Please use [GitHub Issues](https://github.com/zsoftly/zcp-ddns/issues) to report bugs or request features.

When filing a bug report, include:

- The zcp-ddns version
- Your operating system and architecture (or container runtime)
- Your configuration with secrets redacted
- The expected vs. actual behavior
- Any relevant log output

## Pull Requests

We welcome contributions from the community. Before opening a pull request:

1. Claim an issue first by commenting on it. Open a new issue to discuss any change that does not have one. One claimed issue at a time until your first PR is merged.
2. Fork the repository and create a feature branch from `main`.
3. Follow the existing code style by running `make fmt` before committing.
4. Add or update tests for any changed behavior.
5. Sign off every commit (see Developer Certificate of Origin below).
6. Open a pull request with a clear description and reference the issue it closes.

Keep PRs small and focused on one issue. The review target is an initial response within 2 business days.

## Developer Certificate of Origin

This project uses the [DCO](https://developercertificate.org/). Every commit must be signed off:

```bash
git commit -s -m "your message"
```

This adds a `Signed-off-by` line certifying you have the right to submit the contribution under the project license. Maintainers verify the DCO before merge.

## Design decisions

Anything that changes scope, architecture, or acceptance criteria is decided and recorded in the GitHub issue, not in Slack. Slack (`#zcp-open-source` at <https://zcp.zsoftly.ca/community>) is for discussion and getting unblocked.

## Ground rules

- All ZCP API access goes through `github.com/zsoftly/zcp-cli/pkg/api/dns`. If the client is missing something, open an issue in [zcp-cli](https://github.com/zsoftly/zcp-cli) rather than working around it here.
- Never commit secrets, tokens, real customer names, or registrable production domains in code, tests, fixtures, or docs. Use your assigned test zone and placeholder data.
- By contributing you agree your work is licensed under the [Apache License 2.0](LICENSE).
