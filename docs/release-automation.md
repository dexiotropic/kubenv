---
icon: robot
---

# Release automation

This page is for maintainers who need to understand how releases are published.

## What the release flow does

Published releases use:

* **Release Please** for versioning and GitHub releases
* **GoReleaser** for binaries, checksums, and the Argo CD CMP image
* repository scripts for Homebrew and Argo CD release assets
* **krew-release-bot** for public `krew-index` update PRs

## Published artifacts

Each release is intended to publish:

* `kubenv` archives for Linux, macOS, and Windows
* `kubectl-kenv` archives for Linux, macOS, and Windows
* `kubenv-argocd-cmp` archives for Linux
* `checksums.txt`
* a generated Homebrew formula artifact
* a generated Krew manifest asset named `kenv.yaml`
* `ghcr.io/dexiotropic/kubenv-argocd-cmp:<tag>`
* `ghcr.io/dexiotropic/kubenv-argocd-cmp:latest`

## Krew automation

The repository root contains `.krew.yaml`, which is the source of truth for the Krew manifest template.

During releases:

1. the workflow renders `kenv.yaml` from `.krew.yaml`
2. the rendered manifest is uploaded as a release asset
3. `krew-release-bot` opens the public `krew-index` update PR for non-prerelease releases

## Repository settings used by releases

If you maintain releases for this project, these settings matter:

* `RELEASE_PLEASE_TOKEN`
* `HOMEBREW_TAP_REPOSITORY`
* `HOMEBREW_TAP_BRANCH` when needed
* `HOMEBREW_TAP_TOKEN`
