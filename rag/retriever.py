"""Dependency-free hybrid retriever for the washoku seed corpus."""
from __future__ import annotations

import json
import re
from pathlib import Path
from typing import Any

ROOT = Path(__file__).resolve().parent
METADATA = ROOT / "metadata.jsonl"
CHUNKS = ROOT / "chunks.jsonl"


def _read_wrapped_jsonl(path: Path) -> list[dict[str, Any]]:
    """Read the repository's JSONL artifact, including its single JSON wrapper."""
    raw = path.read_text(encoding="utf-8").strip()
    if not raw:
        return []
    try:
        wrapper = json.loads(raw)
        raw = wrapper.get("content", raw)
    except json.JSONDecodeError:
        pass
    return [json.loads(line) for line in raw.splitlines() if line.strip()]


def load_corpus() -> tuple[dict[str, dict[str, Any]], list[dict[str, Any]]]:
    metadata = {item["id"]: item for item in _read_wrapped_jsonl(METADATA)}
    chunks = _read_wrapped_jsonl(CHUNKS)
    return metadata, chunks


def _tokens(text: str) -> set[str]:
    # Japanese is handled as character bigrams; whitespace-separated languages
    # still get useful token matches.
    text = text.lower()
    words = set(re.findall(r"[a-z0-9]+|[ぁ-んァ-ヶ一-龠々ー]", text))
    compact = re.sub(r"\s+", "", text)
    words.update(compact[i : i + 2] for i in range(max(0, len(compact) - 1)))
    return words


def _score(query: str, text: str, metadata: dict[str, Any]) -> float:
    q = _tokens(query)
    if not q:
        return 0.0
    t = _tokens(text)
    lexical = len(q & t) / len(q)

    # Small deterministic boosts for explicit hard/semantic constraints.
    qcompact = re.sub(r"\s+", "", query)
    boost = 0.0
    if "川崎" in qcompact and metadata.get("area") == "川崎":
        boost += 0.25
    if any(x in qcompact for x in ("魚", "海鮮", "刺身", "寿司", "鮨")) and metadata.get("fish") in {"high", "very_high"}:
        boost += 0.15
    if any(x in qcompact for x in ("新鮮", "鮮度", "市場", "直送")) and metadata.get("freshness") == "very_high":
        boost += 0.2
    if any(x in qcompact for x in ("見栄え", "映え", "インスタ", "見た目")) and metadata.get("visual_appeal") == "high":
        boost += 0.15
    if "一人" in qcompact and "solo" in metadata.get("target_use", []):
        boost += 0.1
    return lexical + boost


def search(query: str, *, limit: int = 5, budget_yen_max: int | None = None, area: str | None = None) -> list[dict[str, Any]]:
    metadata, chunks = load_corpus()
    results: list[dict[str, Any]] = []
    for chunk in chunks:
        item = metadata.get(chunk["document_id"], {}).copy()
        if not item:
            continue
        if budget_yen_max is not None and item.get("budget_yen", 10**9) > budget_yen_max:
            continue
        if area and item.get("area") != area:
            continue
        score = _score(query, chunk["text"], item)
        if score <= 0:
            continue
        results.append({"score": round(score, 4), "id": item["id"], "name": item["name"], "text": chunk["text"], "metadata": item})
    return sorted(results, key=lambda x: (-x["score"], x["name"]))[:limit]


def answer(query: str, **kwargs: Any) -> str:
    """Return a grounded answer without requiring an LLM or API key."""
    results = search(query, **kwargs)
    if not results:
        return "該当する候補が見つかりませんでした。"
    lines = [f"質問: {query}", "", "候補:"]
    for i, r in enumerate(results, 1):
        m = r["metadata"]
        lines.append(f"{i}. {r['name']} — 約{m.get('budget_yen')}円 / {m.get('category')} / {r['text']}")
    return "\n".join(lines)
