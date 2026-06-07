# Entrypoints

`kubenv` ships the same renderer through three entrypoints so the behavior stays
consistent across local use, kubectl workflows, and Argo CD.

## Available entrypoints

| Entrypoint | Best for | Page |
| --- | --- | --- |
| `kubenv` | Local rendering or piping directly into `kubectl apply` | [kubenv CLI](KUBENV.md) |
| `kubectl kenv` | kubectl-native usage and Krew installation | [kubectl plugin](KUBECTL.md) |
| `kubenv-argocd-cmp` | GitOps and Argo CD Config Management Plugin workflows | [Argo CD CMP](ARGOCD.md) |

## Next step

Choose the page that matches where rendering should happen:

- [kubenv CLI](KUBENV.md)
- [kubectl plugin](KUBECTL.md)
- [Argo CD CMP](ARGOCD.md)
