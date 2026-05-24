# Argo CD CMP

`kubenv-argocd-cmp` is the Argo CD Config Management Plugin entrypoint.

The example plugin definition lives at:

- `packaging/argocd/plugin.yaml`

Release automation publishes:

- a `ghcr.io/dexiotropic/kubenv-argocd-cmp:<tag>` sidecar image
- a `kubenv-argocd-plugin.yaml` release asset generated from `packaging/argocd/plugin.yaml`

## Which installation model this repo uses

Argo CD CMP still works by adding a **sidecar container** to `argocd-repo-server`.

There are two common ways to do that:

1. use a mostly stock sidecar image and mount `plugin.yaml` from a ConfigMap
2. publish a **custom sidecar image** that already contains the plugin binary and `plugin.yaml`

This repository uses **option 2**.

That means the published artifact is:

- a sidecar image that already contains:
  - `kubenv-argocd-cmp`
  - `plugin.yaml` at `/home/argocd/cmp-server/config/plugin.yaml`
  - the correct sidecar entrypoint pointing at `/var/run/argocd/argocd-cmp-server`

What is **not** automatic is installation into a live Argo CD instance. An operator still needs to patch `argocd-repo-server` (or the Helm/operator values that manage it) to run that sidecar image.

So, in short:

- **publish** = build and publish the ready-to-use sidecar image
- **install** = add that sidecar to `argocd-repo-server`

The stock-image path is still possible, but it is not the model this repository is targeting.

## Sidecar installation requirements

The published image is intentionally small. It does **not** bundle `argocd-cmp-server`; instead it expects Argo CD to provide that binary through the standard shared mount at:

- `/var/run/argocd`

That means a working installation must mount the same core paths Argo CD documents for CMP sidecars:

- `/var/run/argocd`
- `/home/argocd/cmp-server/plugins`
- `/tmp`

The container should also run as:

- `runAsNonRoot: true`
- `runAsUser: 999`

Because `plugin.yaml` is baked into the image, you do **not** need a separate ConfigMap mount for it unless you want to override the baked configuration.

## Example sidecar patch

The exact patch depends on how you manage Argo CD, but the sidecar should look roughly like this:

```yaml
containers:
  - name: kubenv-cmp
    image: ghcr.io/dexiotropic/kubenv-argocd-cmp:v0.1.0
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
```

If your Argo CD installation does not already define `cmp-tmp`, add a separate `emptyDir` volume for it just as the Argo CD documentation recommends.

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

- array parameters can be accessed with '{{ env.VAR_NAME_<index> }}' syntax, for exaple to access the second element of an array parameter named `ITEMS`, you can use `{{ env.ITEMS_1 }}`
- malformed `ARGOCD_APP_PARAMETERS` fails manifest generation
- `ARGOCD_ENV_*` variables are exposed to placeholders without the prefix
- the CMP config currently advertises a `vars` map parameter in `packaging/argocd/plugin.yaml`
