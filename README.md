# washoku

Restaurant knowledge RAG for asking questions over restaurant text and metadata.

## Architecture

```text
Documents
   ↓
Metadata + Chunks
   ↓
RAG retrieval
   ↓
Grounded context
   ↓
Interface
 ┌──────┬──────┬─────┐
 REST   MCP    CLI
```

**RAG is not REST.** RAG is the knowledge/retrieval layer. REST, MCP, and CLI are interchangeable interfaces around it.

## Setup and run

First time:

```bash
git clone https://github.com/bonsai/washoku.git
cd washoku
./set.sh
./run.sh
```

After setup, `run.sh` is runtime-only. It does not clone, pull, or install anything.

One-shot query:

```bash
./run.sh "川崎で2000円以内、新鮮な魚"
```

Search mode:

```text
washoku> search 川崎 魚 新鮮
```

If `./set.sh` or `./run.sh` is not executable on a fresh checkout, run once with:

```bash
bash set.sh
```

`set.sh` creates `.venv`, installs `requirements.txt` when present, runs the retriever tests, and makes both scripts executable.

## Current corpus

Five initial restaurant documents are stored under `documents/`, with structured metadata and retrieval chunks under `rag/`.

The corpus is currently a small, inspectable seed dataset rather than a production vector database.

## Example questions

- 川崎で2000円以内ならどこが一番魚の鮮度が高い？
- 一人で気軽に定食を食べたい
- 2000円で見栄えも重視するとどこ？
- 市場直送・新鮮な魚を優先すると？

See `rag/README.md` and `interface/README.md` for the design.
