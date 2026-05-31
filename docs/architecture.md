# Architecture

If you are evaluating or contributing to the project, the codebase keeps a single renderer and three executables:

- `kubenv` for local CLI use
- `kubectl-kenv` for kubectl plugin use
- `kubenv-argocd-cmp` for Argo CD Config Management Plugin use

All three binaries share the same strict renderer so the behavior stays aligned between local CLI usage, the kubectl plugin, and Argo CD.
