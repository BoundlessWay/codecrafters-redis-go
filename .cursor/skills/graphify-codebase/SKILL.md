---
name: graphify-codebase
description: Rebuilds or interprets the graphify knowledge graph for this Codecrafters Redis-in-Go repo. Use when the user asks to run graphify, refresh or update graphify-out, explain graph outputs, or after meaningful edits under app/, tests/, docs/, or README.md.
---

# Graphify for this repository

## How the agent should use existing outputs

1. Prefer `graphify-out/GRAPH_REPORT.md` for **overview** (communities, god nodes, surprising links, suggested questions).
2. Use `graphify-out/graph.json` for **precise** structure: node labels, edge `relation`, `confidence` (EXTRACTED / INFERRED / AMBIGUOUS), and `source_file`.
3. Open `graphify-out/graph.html` in a browser for **visual** exploration (Cursor has no built-in graph viewer).
4. Project rule `.cursor/rules/graphify-knowledge-graph.mdc` already tells the agent to consult these paths for architecture questions; this skill covers **rebuilding** and **operational** details.

## When to rebuild

- User explicitly asks to graphify / refresh the graph.
- `graphify-out/` is missing or **clearly outdated** after changes to `app/*.go`, `tests/*.go`, `README.md`, or `docs/*.md`.

## Prerequisites

- Python on PATH.
- Package: `pip install graphifyy` (import name `graphify`).
- Run all commands from the **repository root** (`codecrafters-redis-go`).

## Rebuild pipeline (high level)

Execute in order; merge AST + semantic before `build_from_json`.

1. **Detect** — `graphify.detect.detect(Path('.'))` → write JSON summary of file buckets (`code`, `document`, …). On Windows, if the detect file is written from PowerShell, read it with encoding **`utf-8-sig`** to tolerate a BOM.

2. **AST (Go)** — For every path in `detect['files']['code']`, expand dirs with `graphify.extract.collect_files`, then `graphify.extract.extract(code_files)` → `.graphify_ast.json`. This step uses **no LLM tokens**.

3. **Semantic (markdown)** — For `README.md` and `docs/*.md`, produce a fragment with nodes/edges (and optional hyperedges) matching graphify’s extraction schema, then merge into `.graphify_semantic.json`. Prefer **explicit `references` edges** from doc concepts to AST node ids when the doc names a file or symbol. Use **`graphify.cache.save_semantic_cache`** when applicable so reruns skip unchanged files.

4. **Merge** — AST nodes first; add semantic nodes whose `id` is not already present; concatenate edges and hyperedges → `.graphify_extract.json`.

5. **Build and export** — `graphify.build.build_from_json` → `graphify.cluster.cluster` / `score_all` → `graphify.report.generate` → `graphify.export.to_json` and `to_html` under `graphify-out/`. Persist analysis for labeling pass if the workflow uses a second report regeneration.

6. **Post-run** — Update `graphify-out/cost.json` if tracking tokens; remove transient dotfiles (`.graphify_detect.json`, `.graphify_extract.json`, etc.) only **after** cost/manifest steps if the chosen workflow still needs them.

Corpus here is small (~2k words): **token cost is often zero** if semantic work is cached or done without an LLM API; a full upstream graphify “subagent per chunk” run would charge tokens on first extraction.

## Incremental updates

If the project later adopts graphify’s **`--update`** flow, only changed files are re-extracted; merge the new extraction into the existing `graphify-out/graph.json` per graphify docs. After **code-only** changes, AST-only refresh may suffice.

## Optional tooling

- **MCP:** `python -m graphify.serve graphify-out/graph.json` exposes query tools for clients that support MCP (configure in Cursor MCP settings if desired).
- **`.graphify_python`** — optional interpreter path file from some scripts; safe to delete or `.gitignore`; ordinary `python` on PATH is enough.

## Anti-patterns

- Do not claim cross-file relationships that are **not** present in `graph.json` (or clearly EXTRACTED in docs).
- Do not assume Obsidian output unless the user requested `--obsidian`.
