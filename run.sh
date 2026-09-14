#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${ROOT}"

printf '\n=== washoku RAG ===\n\n'

PYTHON="${ROOT}/.venv/bin/python"
if [[ ! -x "${PYTHON}" ]]; then
  echo '環境が未セットアップです。先に ./set.sh を実行してください。' >&2
  exit 1
fi

CLI="${ROOT}/rag/cli.py"
[[ -f "${CLI}" ]] || { echo "rag/cli.py が見つかりません。" >&2; exit 1; }

ask() {
  "${PYTHON}" "${CLI}" ask "$@"
}

search() {
  "${PYTHON}" "${CLI}" search "$@"
}

if [[ $# -gt 0 ]]; then
  ask "$*"
  exit $?
fi

printf 'RAG ready. 質問を入力してください。\n'
printf '終了: exit / 検索: search <query>\n\n'

while true; do
  printf 'washoku> '
  IFS= read -r question || break

  [[ -z "${question}" ]] && continue
  [[ "${question}" == "exit" ]] && break

  if [[ "${question}" == search\ * ]]; then
    search "${question#search }"
  else
    ask "${question}"
  fi
  printf '\n'
done
