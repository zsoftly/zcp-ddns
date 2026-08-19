# Roadmap

This roadmap tracks planned work for `zcp-ddns`. It should stay aligned with GitHub Issues.

## Initial Release

- [x] Public repository scaffolding
- [x] Branch protection and CODEOWNERS
- [x] Minimal buildable Go service entry point
- [x] CI for formatting, vet, tests, and build
- [ ] Configuration loader
- [ ] Public IP detection with fallback providers
- [ ] ZCP DNS API integration through the official ZCP client
- [ ] A record update support
- [ ] AAAA record update support
- [ ] Dry-run mode
- [ ] Run-once mode
- [ ] Docker image publishing
- [ ] Example Docker Compose file

## Later

- [ ] Multiple records per config file
- [ ] Per-record TTL configuration
- [ ] Structured logs
- [ ] Health endpoint
- [ ] Kubernetes deployment example
- [ ] Helm chart or Kustomize example
- [ ] Metrics endpoint

## Out of Scope

- Full DNS zone management.
- Production domain registration.
- Non-DNS ZCP resource management.
- Direct API calls that bypass the official ZCP client.
