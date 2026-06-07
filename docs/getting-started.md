# Getting started

This page gives you the shortest path to a working render flow.

## 1. Pick an entrypoint

Use the entrypoint that matches where rendering should happen:

- `kubenv` for local CLI usage
- `kubectl kenv` for kubectl-native workflows
- `kubenv-argocd-cmp` for Argo CD and GitOps workflows

If you are not sure, start with `kubenv`.

## 2. Choose how to pass variables

Variable precedence is:

1. dotenv files
2. process environment variables
3. command-line arguments such as `--set KEY=VALUE`

That means CLI flags are the final override when you need them.

## 3. Use the default placeholder style

The default syntax is:

```yaml
{{ env.NAME }}
```

Example:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: demo
data:
  message: "{{ env.GREETING }} {{ env.NAME }}"
```

```sh
kubenv render --set GREETING=hello --set NAME=world -f manifest.yaml
```

## 4. Switch modes only when you need to

If you already have manifests written with shell-style placeholders, enable:

```sh
--shell-style
```

That changes parsing to:

```yaml
$NAME
${NAME}
```

If you need to keep a literal explicit placeholder for another tool, escape it:

```yaml
{{ !env.NAME }}
```

That renders back to:

```yaml
{{ env.NAME }}
```

## 5. Continue with the page for your workflow

- [kubenv CLI](KUBENV.md)
- [kubectl plugin](KUBECTL.md)
- [Argo CD CMP](ARGOCD.md)
