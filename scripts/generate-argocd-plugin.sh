#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 || $# -gt 2 ]]; then
  echo "usage: $0 <output-path> [plugin-version]" >&2
  exit 1
fi

output="$1"
plugin_version="${2:-}"

mkdir -p "$(dirname "$output")"

awk -v plugin_version="$plugin_version" '
  {
    print
  }
  $0 == "spec:" && plugin_version != "" {
    print "  version: " plugin_version
  }
' packaging/argocd/plugin.yaml >"${output}"
