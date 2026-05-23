# kubenv

`kubenv` is a minimal manifest renderer for Kubernetes-focused variable substitution. Think of dotenv for Kubernetes manifests. It is designed to be simple, deterministic, and compatible with Argo CD Config Management Plugin (CMP) parameters:

- strict `{{ env.VAR }}` substitution
- deterministic output
- fail on missing variables
- no loops, patches, overlays, or packaging system
- compatible with Argo CD CMP parameters

If you need a full-featured templating engine with loops, conditionals, and a rich function library, consider using tools like Helm or Kustomize. `kubenv` is intentionally minimal to provide a straightforward solution for variable substitution in Kubernetes manifests without the complexity of a full templating system.

The repository keeps the MVP in one place:

- `cmd/kubenv`: main CLI
- `cmd/kubenv-argocd-cmp`: Argo CD Config Management Plugin entrypoint
- `internal/render`: strict substitution engine
- `docs/`: design notes
- `examples/`: sample manifests
- `packaging/krew/`: Krew manifest assets


## Quick start

```sh
go run ./cmd/kubenv render -f examples/configmap.yaml
```

```sh
GREETING=hello go run ./cmd/kubenv render -f examples/configmap.yaml
```

```sh
go run ./cmd/kubenv render --dotenv -f examples/configmap.yaml
```

```sh
go run ./cmd/kubenv render --dotenv-file .env.dev --set NAME=world -f examples/configmap.yaml
```

```sh
go run ./cmd/kubenv render -f first.yaml -f second.yaml
```

```sh
go run ./cmd/kubenv apply --env -f examples/configmap.yaml -- --namespace default
```

**As a Krew plugin:**

Manually build and install:

```sh
go build -o ~/.krew/bin/kubectl-env ./cmd/kubectl-env
```

```sh
kubectl env --dotenv -f examples/configmap.yaml apply --namespace default
```

> [!IMPORTANT]
> Beware of `-f` placement, if you want the file to be rendered with variable substitution, it must come before `apply` and any `kubectl apply` flags. Any file provided after `apply` with the `-f` flag will be passed directly to `kubectl apply` without rendering, which may lead to unexpected results if you intended for it to be rendered first.

```sh

## Variable sources

`kubenv render` loads variables with this precedence:

1. `--set KEY=VALUE`
2. process environment
3. `.env` via `--dotenv` or a specific file via `--dotenv-file`

Notes:

- `--dotenv` loads `.env`
- `--dotenv-file <path>` loads a specific dotenv file
- `--ignore-process-env` disables process environment loading
- `--dotenv` and `--dotenv-file` are mutually exclusive
- `-f` may be repeated and files are rendered in the order provided
- `apply` forwards extra arguments after `--` to `kubectl apply`
- `kubectl kubenv` also supports `kubenv` flags before `apply`, followed by raw `kubectl apply` flags
- dotenv parsing is intentionally minimal: blank lines and `#` comments are supported, and variable lines must be `KEY=VALUE`
