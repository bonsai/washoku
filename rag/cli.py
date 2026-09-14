#!/usr/bin/env python3
"""CLI for local washoku RAG retrieval."""
from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from retriever import answer, search  # noqa: E402


def main() -> int:
    parser = argparse.ArgumentParser(prog="washoku")
    sub = parser.add_subparsers(dest="command", required=True)

    p_search = sub.add_parser("search", help="retrieve grounded restaurant candidates")
    p_search.add_argument("query")
    p_search.add_argument("--limit", type=int, default=5)
    p_search.add_argument("--budget-max", type=int)
    p_search.add_argument("--area")
    p_search.add_argument("--json", action="store_true")

    p_ask = sub.add_parser("ask", help="answer using retrieved corpus only")
    p_ask.add_argument("query")
    p_ask.add_argument("--limit", type=int, default=5)
    p_ask.add_argument("--budget-max", type=int)
    p_ask.add_argument("--area")

    args = parser.parse_args()
    kwargs = {"limit": args.limit, "budget_yen_max": getattr(args, "budget_max", None), "area": getattr(args, "area", None)}
    if args.command == "search":
        results = search(args.query, **kwargs)
        print(json.dumps(results, ensure_ascii=False, indent=2) if args.json else "\n".join(
            f"{i}. {r['name']} [{r['score']}] — {r['text']}" for i, r in enumerate(results, 1)
        ) or "該当する候補が見つかりませんでした。")
    else:
        print(answer(args.query, **kwargs))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
