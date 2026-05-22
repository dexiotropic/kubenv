# kubenv

`kubenv` is a minimal manifest renderer for Kubernetes-focused variable substitution.

The repository keeps the MVP in one place:

- `cmd/kubenv`: main CLI
- `cmd/kubenv-argocd-cmp`: Argo CD Config Management Plugin entrypoint
- `internal/render`: strict substitution engine
- `docs/`: design notes
- `examples/`: sample manifests
- `packaging/krew/`: Krew manifest assets

## MVP goals

- strict `{{ env.VAR }}` substitution
- deterministic output
- fail on missing variables
- no loops, patches, overlays, or packaging system
- compatible with Argo CD CMP parameters

## Quick start

```sh
go run ./cmd/kubenv render -f examples/configmap.yaml
```

```sh
GREETING=hello go run ./cmd/kubenv render -f examples/configmap.yaml
```

```sh
go run ./cmd/kubenv render --env -f examples/configmap.yaml
```

```sh
go run ./cmd/kubenv render --env-file .env.dev --set NAME=world -f examples/configmap.yaml
```

```sh
go run ./cmd/kubenv render -f first.yaml -f second.yaml
```

## Variable sources

`kubenv render` loads variables with this precedence:

1. `--set KEY=VALUE`
2. process environment
3. `.env` via `--env` or a specific file via `--env-file`

Notes:

- `--env` loads `.env`
- `--env-file <path>` loads a specific dotenv file
- `--ignore-process-env` disables process environment loading
- `--env` and `--env-file` are mutually exclusive
- `-f` may be repeated and files are rendered in the order provided
- dotenv parsing is intentionally minimal: blank lines and `#` comments are supported, and variable lines must be `KEY=VALUE`
