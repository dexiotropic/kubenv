# kubenv CLI

Use the `kubenv` binary when you want to render manifests locally with strict `{{ env.NAME }}` substitution.

## Installation

### Brew

```sh
brew tap dexiotropic/homebrew-tap
brew install kubenv
```

### Build from source

```sh
go build -o kubenv ./cmd/kubenv
```

## Help

You can inspect command and flag help directly from the binary:

```sh
kubenv --help
kubenv render --help
kubenv apply --help
```

## Commands

### Render

```sh
kubenv render -f examples/configmap.yaml
kubenv render --dotenv -f examples/configmap.yaml
kubenv render --dotenv-file .env.dev --set NAME=world -f examples/configmap.yaml
kubenv render -f first.yaml -f second.yaml
```

### Apply

```sh
kubenv apply --dotenv -f examples/configmap.yaml -- --namespace default
```

`apply` renders your manifests first and then runs:

```sh
kubectl apply -f -
```

Any arguments after `--` are forwarded to `kubectl apply`.

## Variable sources

`kubenv render` and `kubenv apply` load variables with this precedence:

1. `--set KEY=VALUE`
2. process environment
3. `.env` via `--dotenv` or a specific file via `--dotenv-file`

Keep in mind:

- `--dotenv` loads `.env`
- `--dotenv-file <path>` loads a specific dotenv file
- `--ignore-process-env` disables process environment loading
- `--dotenv` and `--dotenv-file` are mutually exclusive
- `-f` may be repeated and files are rendered in the order provided
- dotenv parsing is intentionally minimal: blank lines and `#` comments are supported, and variable lines must be `KEY=VALUE`
