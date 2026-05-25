# Argo CD CMP

Use `kubenv-argocd-cmp` when you want the same renderer to run inside an Argo CD Config Management Plugin sidecar.

The example plugin definition lives at:

- `packaging/argocd/plugin.yaml`

Release automation publishes:

- a `ghcr.io/dexiotropic/kubenv-argocd-cmp:<tag>` sidecar image
- a moving `ghcr.io/dexiotropic/kubenv-argocd-cmp:latest` sidecar image
- a `kubenv-argocd-plugin.yaml` release asset for the unversioned `latest` plugin name
- a `kubenv-argocd-plugin-versioned.yaml` release asset that matches the published release tag


## Installation

The installation process changes based on your Argo CD management and deployment model, however all these methods involve adding a sidecar container to `argocd-repo-server` that runs the published image. Because `plugin.yaml` is baked into the image, you do **not** need a separate ConfigMap mount for it unless you want to override the baked configuration.

> [!TIP]
> If you prefer automatic sidecar updates, use `:latest` in the next examples. If you prefer reproducible deployments, pin a specific release tag such as `:v0.2.1`. You can find these tags on the [GitHub Container Registry page for this project](https://github.com/dexiotropic/kubenv/pkgs/container/kubenv-argocd-cmp).

> [!IMPORTANT]
> When you use `:latest`, the baked plugin config does **not** set `spec.version`, so your Argo CD applications should use:
> 
> - `spec.source.plugin.name: kubenv`
> 
> When you use a pinned sidecar image such as `:v0.3.0`, the baked plugin config sets `spec.version: v0.3.0`, so your Argo CD applications should use:
> 
> - `spec.source.plugin.name: kubenv-v0.3.0`
>

### Argo CD Operator

If you manage Argo CD with the Operator, you can add the sidecar to your `ArgoCD` custom resource definition:

```yaml
spec:
  # Add the sidecar to the repo-server
  repo:
    sidecarContainers:
      - name: kubenv-cmp
        image: ghcr.io/dexiotropic/kubenv-argocd-cmp:latest
        command: ["/var/run/argocd/argocd-cmp-server"]
        securityContext:
          runAsNonRoot: true
          runAsUser: 999
        volumeMounts:
          - mountPath: /var/run/argocd
            name: var-files
          - mountPath: /home/argocd/cmp-server/plugins
            name: plugins
          - mountPath: /tmp
            name: cmp-tmp
    # Add the plugin tmp volume if not already defined
    volumes:
      - name: cmp-tmp
        emptyDir: {}
```

### Argo CD Helm chart

If you manage Argo CD with the Helm chart, you can add the sidecar to your `values.yaml`:

```yaml
repoServer:
  extraContainers:
    - name: kubenv-cmp
      image: ghcr.io/dexiotropic/kubenv-argocd-cmp:latest
      command: ["/var/run/argocd/argocd-cmp-server"]
      securityContext:
        runAsNonRoot: true
        runAsUser: 999
      volumeMounts:
        - mountPath: /var/run/argocd
          name: var-files
        - mountPath: /home/argocd/cmp-server/plugins
          name: plugins
        - mountPath: /tmp
          name: cmp-tmp
  # Add the plugin tmp volume if not already defined
  volumes:
    - name: cmp-tmp
      emptyDir: {}
```

### Custom installation by patching `argocd-repo-server`

If you manage Argo CD with a custom method, you can patch `argocd-repo-server` directly to add the sidecar. The exact patch depends on your deployment method, but the sidecar container definition should look roughly like the previous examples.

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
      name: kubenv
      parameters:
        - name: GREETING
          string: hello
        - name: NAME
          string: world
```

Or you can pass a single map parameter:

```yaml
spec:
  source:
    plugin:
      name: kubenv
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

## Supported environment variables

You can also set environment variables on the plugin source

```yaml
spec:
  source:
    plugin:
      name: kubenv
      env:
        - name: GREETING
          value: hello
        - name: NAME
          value: world
```

These variables are prefixed with `ARGOCD_ENV_` and the plugin reads them with higher precedence than process environment variables like variables that could have been set with `env` in the sidecar definition or inherited from the repo-server environment.

If you pin the sidecar image to a tagged release, change `spec.source.plugin.name` to match that tag, for example `kubenv-v0.3.0`.

## Notes

- array parameters can be accessed with `{{ env.VAR_NAME_<index> }}` syntax; for example, to access the second element of an array parameter named `ITEMS`, use `{{ env.ITEMS_1 }}`
- malformed `ARGOCD_APP_PARAMETERS` fails manifest generation
- `ARGOCD_ENV_*` variables are exposed to placeholders without the prefix
- the CMP config currently advertises a `vars` map parameter in `packaging/argocd/plugin.yaml`
