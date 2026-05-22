# Architecture

The initial repository keeps a single codebase and two executables:

- `kubenv` for local CLI use
- `kubenv-argocd-cmp` for Argo CD Config Management Plugin use

Both binaries share the same strict renderer so local and Argo CD behavior stay aligned.
