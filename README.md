# washoku

Restaurant knowledge RAG + dependency-free Go CLI.

```text
Documents → RAG update → data/restaurants.json → Go CLI
                         ↑
                    GitHub Actions
```

## Run

```bash
git clone https://github.com/bonsai/washoku.git
cd washoku
go run ./cmd/washoku "川崎で2000円以内、新鮮な魚"
```

Interactive:

```bash
./run.sh
```

No Python, external API, embedding service, or LLM is required.

## Data

- `documents/` — source knowledge
- `rag/` — JSONL RAG source data
- `data/` — generated runtime data
- `cmd/washoku/` — search CLI
- `cmd/rag-update/` — RAG builder
- `.github/workflows/` — growth, AW, release

## Automation

`rag-update.yml` regenerates runtime data in GitHub Actions. `aw.yml` runs a query on demand. `release.yml` builds standalone binaries.
