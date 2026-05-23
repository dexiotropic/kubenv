# Argo CD CMP

`kubenv-argocd-cmp` is the Argo CD Config Management Plugin entrypoint.

The example plugin definition lives at:

- `packaging/argocd/plugin.yaml`

## Variable sources

When kubenv runs inside Argo CD, variable precedence is:

1. CMP plugin parameters from `spec.source.plugin.parameters`
2. user-supplied plugin environment variables exposed as `ARGOCD_ENV_*`
3. the remaining process environment of the CMP sidecar

## Supported parameter styles

You can set individual variables as string parameters:

```yaml
spec:
  source:
    plugin:
      name: kubenv-v0.1.0
      parameters:
        - name: GREETING
          string: hello
        - name: NAME
          string: world
```

Or pass a single map parameter:

```yaml
spec:
  source:
    plugin:
      name: kubenv-v0.1.0
      parameters:
        - name: vars
          map:
            GREETING: hello
            NAME: world
```

Both forms support the same manifest placeholders:

```yaml
data:
  message: "{{ env.GREETING }} {{ env.NAME }}"
```

## Notes

- malformed `ARGOCD_APP_PARAMETERS` fails manifest generation
- `ARGOCD_ENV_*` variables are exposed to placeholders without the prefix
- the CMP config currently advertises a `vars` map parameter in `packaging/argocd/plugin.yaml`
