# Architecture

The repository keeps a single codebase and three executables:

- `kubenv` for local CLI use
- `kubectl-env` for kubectl plugin use
- `kubenv-argocd-cmp` for Argo CD Config Management Plugin use

All three binaries share the same strict renderer so local, kubectl, and Argo CD behavior stay aligned.
