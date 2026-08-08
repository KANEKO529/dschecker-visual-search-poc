# Architecture

## Architecture responsibilities

- `frontend/` handles the PoC UI.
- `backend/` is the Go/Gin API layer and coordinates requests between the frontend, inference service, and existing Rails API.
- `inference/` generates image embeddings using DINOv3 ViT-S.
- The existing Rails application owns product data and PostgreSQL/pgvector storage.
- `item_embeddings` is managed by the Rails application, not directly by the Go backend.
- The Go backend does not directly access the Rails application's database; it communicates with Rails through APIs.

## Backend directory structure

```
backend/
  main.go              # Entrypoint. Builds the router and starts the server.
  internal/
    router/
      router.go         # Creates the gin.Engine and registers routes to handlers.
    handler/
      health.go          # Request handlers, one file per endpoint (e.g. health.go for GET /health).
```

- `main.go` depends on `internal/router` only. It does not know about individual handlers or routes.
- `internal/router` depends on `internal/handler`. It owns route-to-handler wiring and does not contain business logic.
- `internal/handler` has no dependency on `internal/router`. Each file implements the handler(s) for one endpoint.
- `internal/` is not imported by any package outside `backend/`, so this layering is internal to the Go backend module.
- This structure covers entrypoint / routing / handler separation only. Service and repository layers are intentionally not introduced yet, since the Go backend does not access a database directly.