# kubenv documentation

`kubenv` is a small manifest renderer for Kubernetes-focused variable
substitution.

Use it when you want:

- explicit `{{ env.NAME }}` placeholders by default
- fail-fast behavior on missing variables
- one render contract across local CLI, kubectl, and Argo CD
- optional shell-style `$NAME` / `${NAME}` support when needed

## Choose your entrypoint

| Entrypoint | Best for | Page |
| --- | --- | --- |
| `kubenv` | Local rendering or piping directly into `kubectl apply` | [kubenv CLI](KUBENV.md) |
| `kubectl kenv` | kubectl-native usage and Krew installation | [kubectl plugin](KUBECTL.md) |
| `kubenv-argocd-cmp` | GitOps and Argo CD Config Management Plugin workflows | [Argo CD CMP](ARGOCD.md) |

## Quick start

### Install

Choose the installation path that matches how you want to run the renderer:

```sh
# Standalone CLI
brew tap dexiotropic/homebrew-tap
brew install kubenv

# kubectl plugin
kubectl krew install kenv
```

For Argo CD, use the published sidecar image:

```sh
ghcr.io/dexiotropic/kubenv-argocd-cmp:latest
```

### Render your first manifest

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: demo
data:
  message: "{{ env.GREETING }} {{ env.NAME }}"
```

```sh
kubenv render --set GREETING=hello --set NAME=world -f manifest.yaml
```

### Placeholder modes

| Need | Syntax |
| --- | --- |
| Default explicit placeholders | `{{ env.NAME }}` |
| Shell-style compatibility | `$NAME` or `${NAME}` with `--shell-style` |
| Keep a literal explicit placeholder | `{{ !env.NAME }}` |

## Start here

- [Getting started](getting-started.md)
- [Entrypoints](entrypoints.md)
- [kubenv CLI](KUBENV.md)
- [kubectl plugin](KUBECTL.md)
- [Argo CD CMP](ARGOCD.md)
- [Maintainers](maintainers.md)
- [Release automation](release-automation.md)
- [Project architecture](architecture.md)
