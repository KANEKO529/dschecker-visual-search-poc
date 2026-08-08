# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

PoC for DSchecker visual search using DINOv3 and vector similarity search. The project verifies whether DSchecker can identify products using image embeddings instead of OCR-based model number recognition.

## Project documentation

Before making architectural or cross-layer changes, read the relevant documentation under `docs/`.

- `docs/development.md` — development workflow, Issue/PR rules, and testing expectations
- `docs/architecture.md` — system architecture, component responsibilities, data ownership, and integration boundaries


## Repository status

This repository is currently a skeleton: `backend/`, `frontend/`, and `inference/` are empty placeholder directories (each holding only a `.gitkeep`). No build, lint, or test tooling has been added yet. When code is added to one of these directories, check for a README or config file inside it (e.g. `package.json`, `go.mod`, `pyproject.toml`) for the actual commands, and update this file accordingly.

## Tech stack

### This repository

- `frontend/` — Next.js / TypeScript
- `backend/` — Go / Gin
- `inference/` — Python / DINOv3 ViT-S

### Data storage

- PostgreSQL / pgvector are owned by the existing Rails application.
- Product embeddings are stored in the Rails-managed `item_embeddings` table.

### Related Rails application

- The existing DSchecker Rails API is maintained in a separate repository.
- The Rails API is also part of this PoC implementation.
- Some PoC issues explicitly require changes to the Rails application.
- Do not modify or assume Rails-side behavior unless the current issue explicitly targets the Rails application.

## Workflow rules (from docs/development.md)

- One issue = one clear purpose; create a feature branch per issue.
- One issue = one PR, in principle.
- Each issue should document: 目的 (purpose), 現状 (current state), 実装内容 (implementation details), 対象外 (out of scope), 関連仕様 (related specs), 完了条件 (completion criteria), テスト項目 (test items).
- Investigate existing code and the blast radius of a change before implementing.
- Run tests after implementing.
- Don't bundle unrelated refactoring into the same issue.
- Before implementing an issue, work out the implementation approach and confirm which files will change.
- Changes to external systems, including the existing DSchecker Rails API, are out of scope unless explicitly included in the issue.


## Implementation workflow

For each issue:

1. Read the issue and related documentation.
2. Investigate the existing implementation and dependencies.
3. Identify the files likely to be changed and the blast radius.
4. Propose an implementation approach before editing code.
5. Implement only the agreed scope.
6. Run relevant tests and checks.
7. Summarize the changes and test results.

Do not start implementation immediately when asked to work on an issue unless the implementation approach is already clear and agreed upon.

## Design principles

- This repository is a PoC. Prioritize validating the technical hypothesis over production-level completeness.
- Search accuracy is the primary evaluation priority.
- The initial image embedding model is DINOv3 ViT-S.
- Product embeddings are stored separately from items so that one product can have multiple embeddings.
- Avoid changing established architecture or technology choices without first explaining the reason and impact.

## Rails integration

The existing DSchecker Rails API is maintained in a separate repository, but some PoC issues require changes to it.

Planned Rails-related work includes:

- Introduce pgvector to Rails.
- Add the `item_embeddings` table.
- Implement the `Item` / `ItemEmbedding` association.
- Implement ItemEmbedding persistence.
- Implement an internal Embedding registration API.
- Resolve an Item by `modelNumber` and save its Embedding.
- Handle 404 / 400 / 422 / 500 responses.

Rails changes should only be performed when working on an issue whose scope explicitly includes the Rails application.
