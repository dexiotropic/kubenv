---
icon: pen-ruler
---

# Project architecture

`kubenv` keeps one renderer and three executables so behavior stays aligned across local use, kubectl workflows, and Argo CD.

## Executables

The codebase keeps a single renderer and three executables:

* `kubenv` for local CLI use
* `kubectl-kenv` for kubectl plugin use
* `kubenv-argocd-cmp` for Argo CD Config Management Plugin use

## Shared behavior

All three binaries share the same strict renderer so placeholder parsing, variable precedence, and missing-variable behavior stay aligned between local CLI usage, the kubectl plugin, and Argo CD.
