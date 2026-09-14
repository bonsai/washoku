#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${ROOT}"

BIN="${ROOT}/washoku"
if [[ ! -x "${BIN}" ]]; then
  if command -v go >/dev/null 2>&1; then
    go build -o "${BIN}" ./cmd/washoku
  else
    echo 'washoku binaryがありません。Releaseから取得するかGoを入れてください。' >&2
    exit 1
  fi
fi

if [[ $# -gt 0 ]]; then
  exec "${BIN}" "$@"
fi

printf '\n=== washoku ===\n質問を入力してください。終了: exit\n\n'
while true; do
  printf 'washoku> '
  IFS= read -r question || break
  [[ -z "${question}" ]] && continue
  [[ "${question}" == "exit" ]] && break
  "${BIN}" "${question}"
  printf '\n'
done
