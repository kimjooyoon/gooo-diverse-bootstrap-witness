#!/usr/bin/env bash
set -euo pipefail

module='github.com/kimjooyoon/gooo-diverse-bootstrap-witness'
path_a="$(go list -deps ./internal/patha | sort -u)"
path_b="$(go list -deps ./internal/pathb | sort -u)"
shared="$(comm -12 <(printf '%s\n' "$path_a") <(printf '%s\n' "$path_b"))"

unexpected=()
while IFS= read -r dependency; do
  [[ -z "$dependency" ]] && continue
  case "$dependency" in
    "$module/internal/wire") ;;
    "$module/internal"/*) unexpected+=("$dependency") ;;
  esac
done <<< "$shared"

if ((${#unexpected[@]} != 0)); then
  printf 'unexpected shared internal dependencies:\n' >&2
  printf '  %s\n' "${unexpected[@]}" >&2
  exit 1
fi

if cmp -s internal/patha/patha.go internal/pathb/pathb.go; then
  echo 'path implementations must have distinct source files' >&2
  exit 1
fi

if rg -n 'internal/(patha|pathb|verifier|eval|evaluator|lowering)' internal/patha internal/pathb --glob '*.go'; then
  echo 'path-a/path-b must not import each other, the verifier, or an evaluation/lowering package' >&2
  exit 1
fi

echo 'independence audit: CLOSED (shared internal dependency: internal/wire only)'
