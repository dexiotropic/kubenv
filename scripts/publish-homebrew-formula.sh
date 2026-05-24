#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 6 ]]; then
  echo "usage: $0 <formula-path> <formula-name> <tap-repo> <tap-branch> <git-author-name> <git-author-email>" >&2
  exit 1
fi

formula_path="$1"
formula_name="$2"
tap_repo="$3"
tap_branch="$4"
git_author_name="$5"
git_author_email="$6"

if [[ -z "${HOMEBREW_TAP_TOKEN:-}" ]]; then
  echo "HOMEBREW_TAP_TOKEN is required" >&2
  exit 1
fi

workdir="$(mktemp -d)"
trap 'rm -rf "$workdir"' EXIT

git clone --branch "$tap_branch" "https://x-access-token:${HOMEBREW_TAP_TOKEN}@github.com/${tap_repo}.git" "$workdir/tap" >/dev/null 2>&1

mkdir -p "$workdir/tap/Formula"
cp "$formula_path" "$workdir/tap/Formula/${formula_name}.rb"

cd "$workdir/tap"

if ! git status --short -- Formula/"${formula_name}.rb" | grep -q .; then
  echo "No Homebrew formula changes to publish"
  exit 0
fi

git config user.name "$git_author_name"
git config user.email "$git_author_email"
git add Formula/"${formula_name}.rb"
git commit -m "Update ${formula_name} formula"
git push origin "$tap_branch"
