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

## Current corpus

Five initial restaurant documents are stored under `documents/`, with structured metadata and retrieval chunks under `rag/`.

The corpus is currently a small, inspectable seed dataset rather than a production vector database.

## Example questions

- 川崎で2000円以内ならどこが一番魚の鮮度が高い？
- 一人で気軽に定食を食べたい
- 2000円で見栄えも重視するとどこ？
- 市場直送・新鮮な魚を優先すると？

See `rag/README.md` and `interface/README.md` for the design.
