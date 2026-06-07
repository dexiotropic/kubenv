---
icon: house
---

# Overview

`kubenv` is a small manifest renderer for Kubernetes-focused variable substitution.

Use it when you want:

* explicit `{{ env.NAME }}` placeholders by default
* fail-fast behavior on missing variables
* one render contract across local CLI, kubectl, and Argo CD
* optional shell-style `$NAME` / `${NAME}` support when needed

{% hint style="info" %}
Kubenv is a lightweight alternative to use cases that do not require complex rendering capabilities, preventing creating Helm charts or Kustomize jargon to write reusable templates that have simple variable substitution requirements.
{% endhint %}

## Choose your entrypoint

<table data-view="cards"><thead><tr><th align="center"></th><th></th><th data-hidden data-card-target data-type="content-ref"></th><th data-hidden data-card-cover data-type="image">Cover image</th></tr></thead><tbody><tr><td align="center"><strong>kubenv</strong></td><td>Local rendering or piping directly into <code>kubectl apply</code></td><td><a href="KUBENV.md">KUBENV.md</a></td><td><a href=".gitbook/assets/Screenshot_2026-06-07_18-34-59.png">Screenshot_2026-06-07_18-34-59.png</a></td></tr><tr><td align="center"><strong>kubectl kenv</strong></td><td>kubectl-native usage and Krew installation</td><td><a href="KUBECTL.md">KUBECTL.md</a></td><td><a href=".gitbook/assets/Screenshot_2026-06-07_18-35-14.png">Screenshot_2026-06-07_18-35-14.png</a></td></tr><tr><td align="center"><strong>kubenv-argocd-cmp</strong></td><td>GitOps and Argo CD Config Management Plugin workflows</td><td><a href="ARGOCD.md">ARGOCD.md</a></td><td><a href=".gitbook/assets/Screenshot_2026-06-07_18-35-32.png">Screenshot_2026-06-07_18-35-32.png</a></td></tr></tbody></table>

## Quick start

{% hint style="success" icon="gratipay" %}
`kenv` installation with Krew is the fastest way to start
{% endhint %}

### Install

Choose the installation path that matches how you want to run the renderer:

{% tabs %}
{% tab title="kubenv" %}
```shellscript
brew tap dexiotropic/homebrew-tap
brew install kubenv
```
{% endtab %}

{% tab title="kubectl kenv" %}
```shellscript
kubectl krew install kenv
```
{% endtab %}

{% tab title="ArgoCD plugin" %}
```yaml
# Patch argocd repo-server deployment template
spec:
  # Add the sidecar to the repo-server
  containers:
    - name: argocd-repo-server # Do not change the main container
      # ...
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
    # ...
```
{% endtab %}
{% endtabs %}

### Render your first manifest

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: demo
data:
  message: "{{ env.GREETING }} {{ env.NAME }}"
```

{% tabs %}
{% tab title="kubenv" %}
```shellscript
kubenv render --set GREETING=hello --set NAME=world -f manifest.yaml
```
{% endtab %}

{% tab title="kubectl kenv" %}
```shellscript
kubectl kenv render --set GREETING=hello --set NAME=world -f manifest.yaml
```
{% endtab %}
{% endtabs %}

### Placeholder modes

| Need                                                    | Syntax                                    |
| ------------------------------------------------------- | ----------------------------------------- |
| Default explicit placeholders                           | `{{ env.NAME }}`                          |
| Shell-style compatibility                               | `$NAME` or `${NAME}` with `--shell-style` |
| Keep a literal explicit placeholder (skip substitution) | `{{ !env.NAME }}`                         |
