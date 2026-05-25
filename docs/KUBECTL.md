# kubectl env plugin

Use the kubectl plugin when you want the same renderer exposed as:

```sh
kubectl env
```

## Build and install locally

```sh
go build -o ~/.krew/bin/kubectl-env ./cmd/kubectl-env
```

Make sure the output directory is on `PATH`, then verify:

```sh
kubectl plugin list
```

Release automation also generates a Krew manifest release asset named `env.yaml`, which you can use when submitting the plugin to a Krew index repository.

You can inspect plugin help directly:

```sh
kubectl env --help
kubectl env render --help
kubectl env apply --help
```

## Comparison with `kubectl-envsubst`

[`hashmap-kz/kubectl-envsubst`](https://github.com/hashmap-kz/kubectl-envsubst) is the closest existing plugin in this space, but the trade-offs are different.

`kubectl-envsubst` currently supports several apply-oriented features that `kubenv` does not:

- directory inputs with optional recursive traversal
- glob expansion
- remote `-f https://...` manifests
- allow-list and prefix-based filtering over `$VAR` / `${VAR}` placeholders

`kubenv` makes the opposite trade:

- it uses explicit `{{ env.NAME }}` placeholders instead of shell-style expansion
- it shares the same render behavior across `kubenv`, `kubectl env`, and Argo CD CMP
- it supports dotenv files and explicit `--set` overrides in addition to process env

So if you are working with existing `${VAR}` manifests and want behavior close to a stricter `kubectl apply` wrapper, `kubectl-envsubst` may be a better fit today. If you want explicit placeholders and consistent behavior across local and GitOps entrypoints, `kubenv` is the better fit.

## Usage

The main plugin flow is:

```sh
kubectl env --dotenv -f examples/configmap.yaml apply --namespace default
```

This syntax splits arguments like this:

- before `apply`: kubenv flags
- after `apply`: raw `kubectl apply` flags

Direct subcommands also work:

```sh
kubectl env render --dotenv -f examples/configmap.yaml
kubectl env apply --dotenv -f examples/configmap.yaml -- --namespace default
```

> [!IMPORTANT]
> If you want a file rendered by kubenv, place its `-f` flag before `apply`. Any `-f` that appears after `apply` is passed directly to `kubectl apply`.
