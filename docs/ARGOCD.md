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

## Plugin selection

This plugin is configured for explicit use.

Set `spec.source.plugin.name` in your Application and Argo CD will invoke `kubenv` directly without relying on repository discovery or marker files.

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

If part of a manifest must keep a literal explicit placeholder for another tool,
write `{{ !env.NAME }}`. kubenv will render that back to `{{ env.NAME }}`
without substituting it, so you can mix kubenv-rendered fields with literal
placeholders in the same manifest.

To switch the Argo CD CMP to shell-style placeholders for one application, set the `kubenv` map parameter:

```yaml
spec:
  source:
    plugin:
      name: kubenv
      parameters:
        - name: kubenv
          map:
            shell-style: "true"
        - name: vars
          map:
            GREETING: hello
            NAME: world
```

With that option enabled, write placeholders as `$GREETING` or `${NAME}` in the rendered manifests.

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

If you want to set a value from the Argo CD CLI instead of editing the Application manifest directly, use `argocd app set --plugin-env`:

```bash
argocd app set management-applications \
  --config-management-plugin kubenv \
  --plugin-env REPO_URL=https://github.com/your-org/your-repo.git
```

After that, sync the application and use the value in your manifests as `{{ env.REPO_URL }}`.

> [!WARNING]
> If you also set the parameter with `spec.source.plugin.parameters`, the parameter value takes precedence over the environment variable, so `{{ env.REPO_URL }}` would resolve to the parameter value instead of the environment variable value.

If you are using a pinned sidecar image with a versioned plugin name, replace `kubenv` in that command with the matching plugin name such as `kubenv-v0.3.0`.

If you pin the sidecar image to a tagged release, change `spec.source.plugin.name` to match that tag, for example `kubenv-v0.3.0`.

## Notes

- the Argo CD CMP is selected explicitly through `spec.source.plugin.name`; it does not rely on discovery or marker files
- the `kubenv` map parameter currently supports `shell-style: "true"` to switch placeholder parsing to `$NAME` and `${NAME}`
- array parameters can be accessed with `{{ env.VAR_NAME_<index> }}` syntax; for example, to access the second element of an array parameter named `ITEMS`, use `{{ env.ITEMS_1 }}`
- malformed `ARGOCD_APP_PARAMETERS` fails manifest generation
- `ARGOCD_ENV_*` variables are exposed to placeholders without the prefix
- the CMP config currently advertises `vars` and `kubenv` map parameters in `packaging/argocd/plugin.yaml`
