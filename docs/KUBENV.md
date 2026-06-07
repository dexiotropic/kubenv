---
icon: square-terminal
---

# kubenv CLI

Use the `kubenv` binary when you want to render manifests locally with strict `{{ env.NAME }}` substitution.

That explicit placeholder style is the default. If you already have manifests written with shell-style placeholders, you can switch to `$VAR` / `${VAR}` rendering with `--shell-style`. If you need a literal explicit placeholder in rendered output, write `{{ !env.NAME }}` and kubenv will emit `{{ env.NAME }}` without substituting it.

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
kubenv render -f manifests/ --recursive
kubenv render -f 'manifests/*.yaml'
kubenv render -f https://example.com/manifest.yaml
kubenv render --shell-style --set IMAGE_TAG=1.2.3 -f deployment.yaml
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

* `--dotenv` loads `.env`
* `--dotenv-file <path>` loads a specific dotenv file
* `--ignore-process-env` disables process environment loading
* `--dotenv` and `--dotenv-file` are mutually exclusive
* `-f` may be repeated and inputs are rendered in the order provided
* `-f` accepts files, directories, glob patterns, and `http` / `https` URLs
* directories include `.yaml`, `.yml`, and `.json` files
* `--recursive` (or `-R`) walks directory inputs recursively
* `--shell-style` switches placeholder parsing from `{{ env.NAME }}` to `$NAME` / `${NAME}`
* `{{ !env.NAME }}` escapes an explicit placeholder and renders it back as `{{ env.NAME }}`
* dotenv parsing is intentionally minimal: blank lines and `#` comments are supported, and variable lines must be `KEY=VALUE`
