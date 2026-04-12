#!/usr/bin/env python3
"""
postToolUse hook: if the agent edits graph-relevant files, suggest refreshing graphify-out.
Reads JSON from stdin; prints JSON to stdout only (fail-open on errors).
"""
from __future__ import annotations

import json
import sys

CONTEXT = (
    "Graphify: app/, tests/, README.md, or docs/ was edited. "
    "Consider re-running the graphify pipeline from the repo root "
    "(see .cursor/skills/graphify-codebase/SKILL.md) so graphify-out/ stays in sync."
)


def collect_paths(obj: object, out: list[str]) -> None:
    if isinstance(obj, dict):
        for k, v in obj.items():
            if k in ("path", "file_path", "target_file", "file", "relativeWorkspacePath") and isinstance(
                v, str
            ):
                out.append(v)
            collect_paths(v, out)
    elif isinstance(obj, list):
        for item in obj:
            collect_paths(item, out)


def is_graph_relevant(path: str) -> bool:
    p = path.replace("\\", "/").strip()
    pl = p.lower()
    if not pl:
        return False
    if "graphify-out/" in pl or pl.startswith("graphify-out/"):
        return False
    if ".cursor/" in pl or "/.cursor/" in pl:
        return False
    base = pl.rsplit("/", 1)[-1]
    if base == "readme.md":
        return True
    if ("/docs/" in pl or pl.startswith("docs/")) and pl.endswith(".md"):
        return True
    if pl.endswith(".go") and (
        "/app/" in pl or pl.startswith("app/") or "/tests/" in pl or pl.startswith("tests/")
    ):
        return True
    return False


def main() -> int:
    try:
        raw = sys.stdin.read()
        data = json.loads(raw) if raw.strip() else {}
    except json.JSONDecodeError:
        print("{}", flush=True)
        return 0

    paths: list[str] = []
    collect_paths(data, paths)
    for key in ("input", "output", "arguments", "params", "tool_input", "tool_output"):
        if key in data and isinstance(data[key], (dict, list)):
            collect_paths(data[key], paths)

    if not paths or not any(is_graph_relevant(p) for p in paths):
        print("{}", flush=True)
        return 0

    print(json.dumps({"additional_context": CONTEXT}, ensure_ascii=False), flush=True)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
