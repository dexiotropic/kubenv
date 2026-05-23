# kubectl env plugin

The kubectl plugin wraps the same renderer used by `kubenv`, but presents it through:

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

## Usage

Primary plugin flow:

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
