# Architecture

`zcp-ddns` is planned as a single-process Go service that runs inside a container and updates ZSoftly DNS records when the desired address for a workload, device, or site changes. The address may be public or private.

## Runtime Flow

1. Load configuration from a file and environment variables.
2. Authenticate to the ZCP API with a DNS-scoped token.
3. Resolve the desired IPv4 and/or IPv6 value from the configured source.
4. Read the configured DNS record from ZSoftly DNS.
5. Create or update the record only when the resolved value differs.
6. Sleep for the configured interval and repeat.

## Boundaries

The service owns:

- Address resolution for public and private A/AAAA records.
- DNS record comparison.
- DNS record create/update orchestration.
- Logging and safe error handling.
- Container runtime behavior.

The service does not own:

- ZCP account creation.
- DNS zone ownership verification.
- General-purpose DNS management.
- Secret storage beyond reading supplied runtime configuration.
- Broad ZCP platform operations outside DNS.

## Address Sources

Initial address sources should stay simple:

- Public IP lookup, for internet-facing records.
- Local interface address, for private records.
- Static configured value, for environments where another system supplies the address.

Additional sources can be added later if they stay within the DDNS scope.

## ZCP API Client

All ZCP DNS API calls should use the official client from:

`github.com/zsoftly/zcp-cli/pkg/api/dns`

If the client is missing required DNS behavior, update the client upstream in `zcp-cli` first. This keeps API behavior consistent across `zcp-cli`, Terraform, and `zcp-ddns`.

## Failure Handling

The service should prefer safe behavior over noisy updates:

- Do not update DNS when address resolution fails.
- Do not delete records as part of normal operation.
- Retry transient API and network failures with bounded backoff.
- Log enough context to debug the failure, without printing tokens or secrets.
