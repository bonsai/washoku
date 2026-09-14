# Interfaces

RAG and interface are separate layers.

```text
Question
  ↓
Interface
  ↓
RAG retrieval
  ↓
metadata filter + chunk ranking
  ↓
grounded context
```

Supported interface directions:

- REST — HTTP applications
- MCP — agent/tool integration
- CLI — local research and debugging

The retrieval core should not depend on any one transport.
