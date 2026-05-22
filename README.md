# kubenv

`kubenv` is a minimal manifest renderer for Kubernetes-focused variable substitution.

The repository keeps the MVP in one place:

- `cmd/kubenv`: main CLI
- `cmd/kubenv-argocd-cmp`: Argo CD Config Management Plugin entrypoint
- `internal/render`: strict substitution engine
- `docs/`: design notes
- `examples/`: sample manifests
- `packaging/krew/`: Krew manifest assets

## MVP goals

- strict `$VAR` substitution
- deterministic output
- fail on missing variables
- no loops, patches, overlays, or packaging system
- compatible with Argo CD CMP parameters

## Quick start

```sh
go run ./cmd/kubenv render -f examples/configmap.yaml
```

```sh
GREETING=hello go run ./cmd/kubenv render -f examples/configmap.yaml
```
