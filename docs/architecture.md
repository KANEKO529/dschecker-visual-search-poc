# Architecture

## Architecture responsibilities

- `frontend/` handles the PoC UI.
- `backend/` is the Go/Gin API layer and coordinates requests between the frontend, inference service, and existing Rails API.
- `inference/` generates image embeddings using DINOv3 ViT-S.
- The existing Rails application owns product data and PostgreSQL/pgvector storage.
- `item_embeddings` is managed by the Rails application, not directly by the Go backend.
- The Go backend does not directly access the Rails application's database; it communicates with Rails through APIs.