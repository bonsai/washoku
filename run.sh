#!/usr/bin/env bash
set -euo pipefail

REPO="https://github.com/bonsai/washoku.git"
NAME="washoku"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET="${ROOT}/${NAME}"

printf '\n=== washoku RAG ===\n\n'

command -v git >/dev/null 2>&1 || { echo 'git が見つかりません。' >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo 'python3 が見つかりません。' >&2; exit 1; }

if [[ ! -d "${TARGET}/.git" ]]; then
  echo "Clone: ${REPO}"
  git clone "${REPO}" "${TARGET}"
else
  echo "Update: ${TARGET}"
  git -C "${TARGET}" pull --ff-only
fi

CLI="${TARGET}/rag/cli.py"
[[ -f "${CLI}" ]] || { echo "rag/cli.py が見つかりません。" >&2; exit 1; }

ask() {
  python3 "${CLI}" ask "$@"
}

search() {
  python3 "${CLI}" search "$@"
}

if [[ $# -gt 0 ]]; then
  ask "$*"
  exit $?
fi

printf '\nRAG ready. 質問を入力してください。\n'
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
