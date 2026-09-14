# Washoku RAG

`washoku` is a small restaurant knowledge base designed for question answering over structured restaurant metadata and source text.

## Layers

1. `documents/` — human-readable source documents.
2. `rag/metadata.jsonl` — one metadata record per restaurant.
3. `rag/chunks.jsonl` — retrieval units with metadata attached.
4. `rag/retriever.py` — dependency-free hybrid retrieval and grounded answer generation.
5. `rag/cli.py` — local CLI interface.
6. `interface/` — protocols used to expose retrieval without coupling RAG to REST.

## Retrieval model

Use hybrid retrieval:

- metadata filtering first when the question contains hard constraints such as area, budget, seating, or meal
- lexical retrieval over chunk text
- deterministic boosts for metadata fit such as fish, freshness, visual appeal, and solo use
- return document IDs and source text so an answer can be grounded

The current implementation intentionally has **no external API and no required dependency**. It is a small executable baseline before adding embeddings, a vector database, or an LLM answer layer.

## Run locally

From `rag/`:

```bash
python cli.py search "川崎で2000円以内、新鮮な魚"
python cli.py search "川崎 魚" --budget-max 2000 --json
python cli.py ask "川崎で一人で魚を食べたい"
```

Tests:

```bash
python -m pytest test_retriever.py
```

## Architecture

```text
question
   ↓
metadata filters
   ↓
lexical retrieval
   ↓
metadata-fit ranking
   ↓
grounded candidates
   ↓
(optional LLM)
```

RAG is the knowledge/retrieval layer. REST, MCP, and CLI are interfaces to that layer, not definitions of RAG itself.
