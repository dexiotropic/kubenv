#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 4 ]]; then
  echo "usage: $0 <tag> <version> <github-repository> <output-path>" >&2
  exit 1
fi

tag="$1"
version="$2"
repository="$3"
output="$4"
dist_dir="$(dirname "$output")"

mkdir -p "$(dirname "$output")"

platform_block() {
  local os="$1"
  local arch="$2"
  local ext="$3"
  local bin="kubectl-env"

  if [[ "$os" == "windows" ]]; then
    bin="kubectl-env.exe"
  fi

  local artifact="kubenv_kubectl-env_${version}_${os}_${arch}.${ext}"
  local sha
  sha="$(sha256sum "${dist_dir}/${artifact}" | awk '{print $1}')"

  cat <<EOF
  - selector:
      matchLabels:
        os: ${os}
        arch: ${arch}
    uri: https://github.com/${repository}/releases/download/${tag}/${artifact}
    sha256: ${sha}
    bin: ${bin}
EOF
}

platforms=(
  "linux amd64 tar.gz"
  "linux arm64 tar.gz"
  "darwin amd64 tar.gz"
  "darwin arm64 tar.gz"
  "windows amd64 zip"
  "windows arm64 zip"
)

{
  cat <<EOF
apiVersion: krew.googlecontainertools.github.com/v1alpha2
kind: Plugin
metadata:
  name: env
spec:
  version: "${tag}"
  homepage: https://github.com/dexiotropic/kubenv
  shortDescription: Render manifests before kubectl apply
  description: |
    kubenv renders Kubernetes manifests using strict variable substitution and
    exposes the renderer as the kubectl env plugin.
  platforms:
EOF

  for platform in "${platforms[@]}"; do
    # shellcheck disable=SC2086
    platform_block ${platform}
  done
} >"${output}"
