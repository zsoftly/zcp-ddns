# Configuration

Configuration is planned through a small YAML file and environment variables.

This document describes the intended shape. The implementation is still in early development.

## Planned YAML Example

```yaml
api_url: https://api.zcp.zsoftly.ca/api
interval: 300s

records:
  - zone: example.com
    name: home.example.com
    type: A
    ttl: 300
```

## Planned Environment Variables

| Variable          | Purpose                                      |
| ----------------- | -------------------------------------------- |
| `ZCP_API_URL`     | Override the ZCP API URL                     |
| `ZCP_TOKEN`       | DNS-scoped ZCP API token                     |
| `ZCP_DDNS_CONFIG` | Path to the YAML configuration file          |
| `ZCP_DDNS_ONCE`   | Run one update cycle and exit, useful for CI |

## Token Scope

Use a DNS-scoped token only. The service should not require broad ZCP account permissions.

## Secrets

Do not store tokens in the repository. For Docker deployments, pass secrets through your container runtime, orchestrator secret store, or environment injection mechanism.

## Record Types

Initial support is planned for:

- `A`
- `AAAA`

Other record types are outside the initial DDNS scope.
