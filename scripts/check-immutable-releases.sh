#!/usr/bin/env bash
set -euo pipefail

repository="${1:-${GITHUB_REPOSITORY:-}}"
api_root="${GITHUB_API_URL:-https://api.github.com}"
token="${GITHUB_TOKEN:-}"

if [[ ! "$repository" =~ ^[^/]+/[^/]+$ ]]; then
  echo 'immutable-release guard: repository is missing or malformed' >&2
  exit 1
fi
if [[ -z "$token" ]]; then
  echo 'immutable-release guard: GITHUB_TOKEN is missing' >&2
  exit 1
fi

payload="$(curl --fail --silent --show-error \
  -H 'Accept: application/vnd.github+json' \
  -H 'X-GitHub-Api-Version: 2026-03-10' \
  -H "Authorization: Bearer $token" \
  "$api_root/repos/$repository/immutable-releases")"

if ! jq -e '.enabled == true' >/dev/null <<<"$payload"; then
  echo 'immutable-release guard: repository immutable releases are not enabled' >&2
  exit 1
fi

echo 'immutable-release guard: CLOSED (enabled=true)'
