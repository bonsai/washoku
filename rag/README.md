# Washoku RAG

`washoku` is a small restaurant knowledge base designed for question answering over structured restaurant metadata and source text.

## Layers

1. `documents/` — human-readable source documents.
2. `rag/metadata.jsonl` — one metadata record per restaurant.
3. `rag/chunks.jsonl` — retrieval units with metadata attached.
4. `interface/` — protocols used to expose retrieval without coupling RAG to REST.

## Retrieval model

Use hybrid retrieval:

- metadata filtering first when the question contains hard constraints such as area, budget, seating, or meal
- lexical/semantic retrieval over chunk text
- rank by relevance plus metadata fit
- return document IDs and source text so an answer can be grounded

RAG is the knowledge/retrieval layer. REST, MCP, and CLI are interfaces to that layer, not definitions of RAG itself.
