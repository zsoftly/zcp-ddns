# ZCP Dynamic DNS (zcp-ddns)

A lightweight dynamic DNS service for the ZSoftly Cloud Platform.

---

## Overview

`zcp-ddns` is a small Go service planned for delivery as a Docker container. It will keep a ZSoftly DNS record pointed at a machine whose public IP address changes. It will authenticate with the ZCP API, periodically detect the current public IP, and create or update configured DNS records automatically.

Typical use cases:

- Home labs and offices on residential or dynamic-IP connections
- Self-hosted services that need a stable hostname
- Edge devices and remote sites without a static IP allocation

## Status

**Early development.** This repository currently contains the project scaffolding only. Implementation is tracked in [GitHub Issues](https://github.com/zsoftly/zcp-ddns/issues) and is open to community contributors. See [ONBOARDING.md](ONBOARDING.md) to get started.

## How it will work

1. Authenticate against the ZCP API with a DNS-scoped API token.
2. Detect the current public IP using multiple public IP detection services with fallback.
3. Compare against the configured DNS records in ZSoftly DNS.
4. Create or update A/AAAA records when the IP changes, then repeat on an interval.

Configuration is planned through a config file and environment variables. The service will run as a single container and will rely on the ZCP API and public IP detection services at runtime.

## Design principles

- **Reuse the official API client.** All ZCP API access goes through [`github.com/zsoftly/zcp-cli/pkg/api/dns`](https://github.com/zsoftly/zcp-cli). Gaps in the client are fixed upstream in zcp-cli, not worked around here.
- **Small and boring.** One binary, one container, minimal configuration, sensible defaults.
- **Safe by default.** DNS-scoped tokens only. The service never needs broad platform permissions.

## Related projects

| Project                                                                     | Purpose                                                                |
| --------------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| [zcp-cli](https://github.com/zsoftly/zcp-cli)                               | Official ZCP command-line interface and Go API client (`pkg/api/dns`)  |
| [terraform-provider-zcp](https://github.com/zsoftly/terraform-provider-zcp) | Terraform/OpenTofu provider for ZCP, including DNS domains and records |
| [zcp-docs](https://docs.zcp.zsoftly.ca)                                     | ZSoftly Cloud Platform documentation                                   |

## Resources

- ZCP DNS documentation: <https://docs.zcp.zsoftly.ca>
- ZCP API base: <https://api.zcp.zsoftly.ca/api>
- Community Slack: <https://zcp.zsoftly.ca/community> (channel `#zcp-open-source`)
- Sign up for a ZCP account: <https://cloud.zcp.zsoftly.ca/register>

## Contributing

Contributions are welcome. Start with [ONBOARDING.md](ONBOARDING.md) for account setup, test credentials, and your first issue, then read [CONTRIBUTING.md](CONTRIBUTING.md) for the pull request process. Maintainership is described in [MAINTAINERS.md](MAINTAINERS.md).

## License

[Apache License 2.0](LICENSE). See [NOTICE](NOTICE).
