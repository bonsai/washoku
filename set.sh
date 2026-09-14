#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${ROOT}"

printf '\n=== washoku setup ===\n\n'
command -v go >/dev/null 2>&1 || { echo 'Goが見つかりません。' >&2; exit 1; }
go version
go build -o washoku ./cmd/washoku
chmod +x set.sh run.sh washoku

printf '\nSetup complete. Run: ./run.sh\n\n'
