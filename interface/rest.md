# REST Interface

REST is an optional transport for the Washoku RAG core.

Suggested endpoints:

- `POST /search` — retrieve relevant restaurant chunks
- `POST /ask` — retrieve context and generate a grounded answer
- `GET /documents/{id}` — inspect one source document

Example search request:

```json
{
  "query": "川崎で2000円以内、新鮮な魚、座って食べたい",
  "filters": {
    "area": "川崎",
    "budget_yen_max": 2000,
    "seating": "seated"
  }
}
```

The API is deliberately transport-specific; retrieval semantics remain in the RAG layer.
