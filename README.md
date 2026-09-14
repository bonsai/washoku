# washoku

Restaurant knowledge RAG with a dependency-free Go CLI.

## Architecture

```text
Documents
   ↓
RAG update (GitHub Actions)
   ↓
data/restaurants.json
   ↓
Go CLI
   ↓
GitHub Release
```

RAG is the knowledge layer. The CLI is the local execution layer. RAG data is regenerated in the cloud; the released CLI only reads the generated data.

## Setup

```bash
git clone https://github.com/bonsai/washoku.git
cd washoku
./set.sh
```

No Python, virtualenv, or API key is required.

## Run

```bash
./run.sh "川崎で2000円以内、新鮮な魚"
```

Or interactive mode:

```bash
./run.sh
```

## RAG update

Changes under `documents/` or the RAG source artifacts trigger `.github/workflows/rag-update.yml`. The workflow runs the dependency-free Go updater in the cloud and commits the generated `data/restaurants.json`.

The local CLI does not crawl, embed, call an LLM, or require an API.

## Release

Tag a version and GitHub Actions builds binaries for Linux amd64, Windows amd64, and macOS arm64:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release assets are standalone Go binaries.
