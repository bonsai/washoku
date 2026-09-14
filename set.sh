#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${ROOT}"

printf '\n=== washoku setup ===\n\n'

command -v python3 >/dev/null 2>&1 || {
  echo 'python3 が見つかりません。' >&2
  exit 1
}

if [[ ! -d .venv ]]; then
  echo 'Creating .venv ...'
  python3 -m venv .venv
fi

PYTHON="${ROOT}/.venv/bin/python"
"${PYTHON}" -m pip install --upgrade pip

if [[ -f requirements.txt ]]; then
  "${PYTHON}" -m pip install -r requirements.txt
fi

if [[ -f rag/test_retriever.py ]]; then
  echo 'Running tests ...'
  "${PYTHON}" rag/test_retriever.py
fi

chmod +x set.sh run.sh

printf '\nSetup complete. Run: ./run.sh\n\n'
