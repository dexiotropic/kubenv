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
- `packaging/krew/`: Krew manifest templates
- `.goreleaser.yaml`: multi-platform release configuration
- `.github/workflows/`: CI and release automation

## Start here

For the direct CLI:

- [`docs/KUBENV.md`](docs/KUBENV.md)

For the kubectl plugin:

- [`docs/KUBECTL.md`](docs/KUBECTL.md)

For Argo CD CMP integration:

- [`docs/ARGOCD.md`](docs/ARGOCD.md)

For internal structure notes:

- [`docs/architecture.md`](docs/architecture.md)

## Release automation

The repository is set up to use:

- **Release Please** for conventional-commit-driven versioning and GitHub releases
- **GoReleaser** for multi-platform binaries, checksums, and the Argo CD CMP image
- repo scripts for generating the Homebrew formula and Krew manifest release assets

Release outputs are intended to include:

- `kubenv` archives for Linux, macOS, and Windows
- `kubectl-env` archives for Linux, macOS, and Windows
- `kubenv-argocd-cmp` archives for Linux
- `checksums.txt`
- a generated Homebrew formula artifact
- a generated Krew plugin manifest
- a published `ghcr.io/dexiotropic/kubenv-argocd-cmp:<tag>` image

If you also want the Homebrew tap updated automatically after each published release, set:

- repository variable `HOMEBREW_TAP_REPOSITORY` to the tap repo, for example `dexiotropic/homebrew-tap`
- optional repository variable `HOMEBREW_TAP_BRANCH` if it is not `main`
- repository secret `HOMEBREW_TAP_TOKEN` with write access to that tap repository
