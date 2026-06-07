# kubectl kenv plugin

Use the kubectl plugin when you want the same renderer exposed as:

```sh
kubectl kenv
```

## Build and install locally

```sh
go build -o ~/.krew/bin/kubectl-kenv ./cmd/kubectl-kenv
```

Make sure the output directory is on `PATH`, then verify:

```sh
kubectl plugin list
```

Release automation generates a `kenv.yaml` Krew manifest asset from `.krew.yaml`
and also uses `krew-release-bot` to open public `krew-index` update PRs after
each non-prerelease GitHub release.

Once the plugin is available in the public Krew index, install it with:

```sh
kubectl krew install kenv
```

The default placeholder style is still `{{ env.NAME }}`. If you need compatibility with existing shell-style manifests, add `--shell-style` before `apply` or before the direct `render` / `apply` subcommand. If you need a literal explicit placeholder in output, write `{{ !env.NAME }}` and kubenv will emit `{{ env.NAME }}` without substituting it.

You can inspect plugin help directly:

```sh
kubectl kenv --help
kubectl kenv render --help
kubectl kenv apply --help
```

## Comparison with `kubectl-envsubst`

[`hashmap-kz/kubectl-envsubst`](https://github.com/hashmap-kz/kubectl-envsubst) is the closest existing plugin in this space, but the trade-offs are different.

`kubectl-envsubst` still supports one notable capability that `kubenv` does not:

- allow-list and prefix-based filtering over `$VAR` / `${VAR}` placeholders

`kubenv` makes the opposite trade:

- it uses explicit `{{ env.NAME }}` placeholders instead of shell-style expansion
- it shares the same render behavior across `kubenv`, `kubectl kenv`, and Argo CD CMP
- it supports dotenv files and explicit `--set` overrides in addition to process env
- it supports files, directories, glob patterns, recursive directory traversal, stdin, and remote `-f https://...` manifests
- it can also opt into `$VAR` / `${VAR}` rendering with `--shell-style` when you need that compatibility

So if you are working with existing `${VAR}` manifests and want behavior close to a stricter `kubectl apply` wrapper, `kubectl-envsubst` may be a better fit today. If you want explicit placeholders and consistent behavior across local and GitOps entrypoints, `kubenv` is the better fit.

## Usage

The main plugin flow is:

```sh
kubectl kenv --dotenv -f examples/configmap.yaml apply --namespace default
kubectl kenv -f manifests/ --recursive apply --namespace default
kubectl kenv -f 'manifests/*.yaml' apply --dry-run=client -o yaml
kubectl kenv -f https://example.com/manifest.yaml apply --namespace default
kubectl kenv --shell-style --set IMAGE_TAG=1.2.3 -f deployment.yaml apply --namespace default
```

This syntax splits arguments like this:

- before `apply`: kubenv flags
- after `apply`: raw `kubectl apply` flags

Direct subcommands also work:

```sh
kubectl kenv render --dotenv -f examples/configmap.yaml
kubectl kenv apply --dotenv -f examples/configmap.yaml -- --namespace default
kubectl kenv render -f manifests/ --recursive
kubectl kenv render --shell-style -f deployment.yaml
```

> [!IMPORTANT]
> If you want a file rendered by kubenv, place its `-f` flag before `apply`. Any `-f` that appears after `apply` is passed directly to `kubectl apply`.
