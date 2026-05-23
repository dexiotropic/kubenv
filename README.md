# kubenv

`kubenv` is a minimal manifest renderer for Kubernetes-focused variable substitution.

It is intentionally narrow in scope:

- strict `{{ env.VAR }}` substitution
- deterministic output
- fail on missing variables
- no loops, overlays, or packaging model
- compatible with Argo CD Config Management Plugin parameters

If you need a full templating system with conditionals, loops, or chart packaging, Helm and Kustomize are a better fit. `kubenv` is aimed at the smaller problem of making a few manifest variables explicit and safe.

## Repository overview

This repository currently ships three entrypoints that share the same renderer:

| Entrypoint | Purpose | Docs |
| --- | --- | --- |
| `kubenv` | Direct CLI for rendering or applying manifests | [`docs/KUBENV.md`](docs/KUBENV.md) |
| `kubectl env` | kubectl plugin wrapper around the same renderer | [`docs/KUBECTL.md`](docs/KUBECTL.md) |
| `kubenv-argocd-cmp` | Argo CD Config Management Plugin entrypoint | [`docs/ARGOCD.md`](docs/ARGOCD.md) |

Implementation layout:

- `cmd/kubenv`: main CLI binary
- `cmd/kubectl-env`: kubectl plugin binary
- `cmd/kubenv-argocd-cmp`: Argo CD CMP binary
- `internal/render`: placeholder parsing and substitution
- `internal/cli`: direct CLI and kubectl plugin orchestration
- `internal/cmp`: Argo CD CMP-specific parameter handling
- `docs/`: focused entrypoint documentation
- `examples/`: sample manifests
- `packaging/argocd/`: CMP configuration assets

## Start here

For the direct CLI:

- [`docs/KUBENV.md`](docs/KUBENV.md)

For the kubectl plugin:

- [`docs/KUBECTL.md`](docs/KUBECTL.md)

For Argo CD CMP integration:

- [`docs/ARGOCD.md`](docs/ARGOCD.md)

For internal structure notes:

- [`docs/architecture.md`](docs/architecture.md)
