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

artifact_sha() {
  sha256sum "${dist_dir}/$1" | awk '{print $1}'
}

cat >"${output}" <<EOF
class Kubenv < Formula
  desc "Strict variable substitution for Kubernetes manifests"
  homepage "https://github.com/dexiotropic/kubenv"
  version "${version}"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/${repository}/releases/download/${tag}/kubenv_kubenv_${version}_darwin_arm64.tar.gz"
      sha256 "$(artifact_sha "kubenv_kubenv_${version}_darwin_arm64.tar.gz")"
    else
      url "https://github.com/${repository}/releases/download/${tag}/kubenv_kubenv_${version}_darwin_amd64.tar.gz"
      sha256 "$(artifact_sha "kubenv_kubenv_${version}_darwin_amd64.tar.gz")"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/${repository}/releases/download/${tag}/kubenv_kubenv_${version}_linux_arm64.tar.gz"
      sha256 "$(artifact_sha "kubenv_kubenv_${version}_linux_arm64.tar.gz")"
    else
      url "https://github.com/${repository}/releases/download/${tag}/kubenv_kubenv_${version}_linux_amd64.tar.gz"
      sha256 "$(artifact_sha "kubenv_kubenv_${version}_linux_amd64.tar.gz")"
    end
  end

  def install
    bin.install "kubenv"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/kubenv version")
  end
end
EOF
