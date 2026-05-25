# kubenv

[![CI](https://github.com/dexiotropic/kubenv/actions/workflows/ci.yml/badge.svg?branch=main&event=push)](https://github.com/dexiotropic/kubenv/actions/workflows/ci.yml)
[![Coverage](https://github.com/dexiotropic/kubenv/actions/workflows/coverage.yml/badge.svg?branch=main&event=push)](https://github.com/dexiotropic/kubenv/actions/workflows/coverage.yml)

`kubenv` is a minimal manifest renderer for Kubernetes-focused variable substitution.

It is intentionally narrow in scope:

- strict `{{ env.VAR }}` substitution
- deterministic output
- fail on missing variables
- no loops, overlays, or packaging model
- compatible with Argo CD Config Management Plugin parameters

Use `kubenv` when you want a small, explicit renderer for a few manifest variables without moving into Helm or Kustomize territory. If you need conditionals, loops, overlays, or chart packaging, Helm and Kustomize are a better fit.

The coverage badge tracks the repository's GitHub Actions coverage gate, which currently requires at least 75% total Go statement coverage.

Compared to a generic `envsubst` wrapper, `kubenv` is intentionally stricter:

- explicit `{{ env.NAME }}` placeholders instead of shell-style expansion
- fail-fast behavior on missing variables
- the same renderer exposed through the direct CLI, `kubectl env`, and Argo CD CMP
- explicit dotenv and `--set` inputs instead of relying only on ambient shell state

If you need shell-style placeholders for an existing manifest set, you can opt into `$VAR` / `${VAR}` rendering in the CLI and kubectl plugin with `--shell-style` while keeping `{{ env.NAME }}` as the default mode.

## Comparison with `kubectl-envsubst`

[`hashmap-kz/kubectl-envsubst`](https://github.com/hashmap-kz/kubectl-envsubst) is a good reference point for this space, and it solves a slightly different problem.

| Area | `kubectl-envsubst` | `kubenv` |
| --- | --- | --- |
| Placeholder syntax | Shell-style `$VAR` / `${VAR}` | Explicit `{{ env.NAME }}` by default, optional `--shell-style` for CLI and kubectl plugin |
| Main interface | `kubectl envsubst apply ...` | `kubenv`, `kubectl env`, and Argo CD CMP |
| Variable source model | Process env filtered by allowed vars or prefixes | `--set`, process env, dotenv, and CMP parameters |
| Safety model | Allow-list and prefix filters reduce accidental substitutions | Placeholder syntax avoids shell-style collisions and missing variables fail immediately |
| File handling | Supports directory walking, glob expansion, recursive mode, stdin, and remote URLs | Supports files, directories, glob expansion, recursive mode, stdin, and remote URLs |

So the difference is not "better at everything"; it is "safer and more explicit for a narrower workflow."

`kubectl-envsubst` is still ahead in one important area for `${VAR}`-style workflows:

- allow-list or prefix-based substitution controls for shell-style placeholders

`kubenv` is stronger when you want:

- placeholders that do not collide with `$`-heavy configs such as shell snippets or NGINX config
- the same render contract across local CLI usage, `kubectl`, and Argo CD
- explicit non-shell inputs such as dotenv files, `--set`, and Argo CD plugin parameters
- deterministic rendering from repeated `-f` inputs, directories, globs, and remote URLs

If you want a near-drop-in `kubectl apply` preprocessor for existing `${VAR}` manifests with optional allow-list or prefix filtering, `kubectl-envsubst` may still be the better fit today. If you want a small renderer with explicit placeholders by default, optional shell-style compatibility, and a cleaner GitOps and Argo CD story, `kubenv` is a strong fit.

## Choose your entrypoint

You can use the same renderer through three entrypoints:

| Entrypoint | Purpose | Docs |
| --- | --- | --- |
| `kubenv` | Direct CLI for rendering or applying manifests | [`docs/KUBENV.md`](docs/KUBENV.md) |
| `kubectl env` | kubectl plugin wrapper around the same renderer | [`docs/KUBECTL.md`](docs/KUBECTL.md) |
| `kubenv-argocd-cmp` | Argo CD Config Management Plugin entrypoint | [`docs/ARGOCD.md`](docs/ARGOCD.md) |

If you are looking around the source tree:

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

If you want the direct CLI:

- [`docs/KUBENV.md`](docs/KUBENV.md)

If you want the kubectl plugin:

- [`docs/KUBECTL.md`](docs/KUBECTL.md)

If you want Argo CD CMP integration:

- [`docs/ARGOCD.md`](docs/ARGOCD.md)

If you want a quick codebase overview:

- [`docs/architecture.md`](docs/architecture.md)
- [`CONTRIBUTING.md`](CONTRIBUTING.md)

## Release automation

If you are maintaining releases for this project, the release flow uses:

- **Release Please** for conventional-commit-driven versioning and GitHub releases
- **GoReleaser** for multi-platform binaries, checksums, and the Argo CD CMP image
- repo scripts for generating the Homebrew formula and Krew manifest release assets

Published releases are intended to include:

- `kubenv` archives for Linux, macOS, and Windows
- `kubectl-env` archives for Linux, macOS, and Windows
- `kubenv-argocd-cmp` archives for Linux
- `checksums.txt`
- a generated Homebrew formula artifact
- a generated Krew plugin manifest
- published `ghcr.io/dexiotropic/kubenv-argocd-cmp:<tag>` and `ghcr.io/dexiotropic/kubenv-argocd-cmp:latest` images

Use the generated `env.yaml` release asset when you open a `krew-index` submission PR.

If you also want the Homebrew tap updated automatically after each published release, set:

- repository variable `HOMEBREW_TAP_REPOSITORY` to the tap repo, for example `dexiotropic/homebrew-tap`
- optional repository variable `HOMEBREW_TAP_BRANCH` if it is not `main`
- repository secret `HOMEBREW_TAP_TOKEN` with write access to that tap repository

Release Please uses an explicit repository secret:

- `RELEASE_PLEASE_TOKEN`

That token needs permission to create branches, commits, tags, releases, and pull requests in the `kubenv` repository. This avoids depending on the repository-level **“Allow GitHub Actions to create and approve pull requests”** setting for the default `GITHUB_TOKEN`.
