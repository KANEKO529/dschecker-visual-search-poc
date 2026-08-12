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
  main.go              # Entrypoint. Loads environment variables, builds the router, and starts the server.
  internal/
    router/
      router.go         # Creates the gin.Engine and registers routes to handlers.
    handler/
      health.go          # Request handlers, one file per endpoint (e.g. health.go for GET /health).
    service/
      embedding_registration.go  # Orchestrates flows that call more than one internal/client in sequence.
    client/
      inference.go       # Outbound HTTP client to the Python inference service.
      rails.go            # Outbound HTTP client to the Rails internal API.
```

- `main.go` depends on `internal/router` only. It does not know about individual handlers or routes.
- `internal/router` depends on `internal/handler`. It owns route-to-handler wiring and does not contain business logic.
- `internal/handler` has no dependency on `internal/router`. Each file implements the handler(s) for one endpoint. A handler depends on `internal/client` directly when it only needs a single external call, or on `internal/service` when the endpoint orchestrates multiple external calls in sequence.
- `internal/service` orchestrates flows that call more than one `internal/client` in sequence (e.g. generating an embedding via the inference service, then registering it with Rails). It has no dependency on `internal/handler` or `internal/router`, and does not depend on `gin`. It does not access a database.
- `internal/client` is responsible for outbound HTTP communication with external services (the inference service and the Rails internal API). It has no dependency on `internal/handler`, `internal/service`, or `internal/router`, and does not depend on `gin`.
- `internal/` is not imported by any package outside `backend/`, so this layering is internal to the Go backend module.
- This structure covers entrypoint / routing / handler / service / outbound-client separation. A repository layer is intentionally not introduced yet, since the Go backend does not access a database directly.

## Inference directory structure

```
inference/
  main.py              # FastAPI app and routes. Contains no embedding generation logic.
  embedding/
    model.py             # Loads DINOv3 ViT-S and generates image embeddings.
```

- `embedding/model.py` has no dependency on FastAPI or `main.py`. It exposes a plain Python function that can be called independently of the HTTP layer.
- The model and image processor are loaded once (lazily, on first use) and reused across calls, not reloaded per request.
- `main.py` does not depend on `embedding/` yet; wiring an endpoint to the embedding generation function is left to a future Issue.

## Frontend directory structure

```
frontend/
  app/                        # Pages and routing only (Next.js App Router convention).
  features/
    <feature-name>/            # One directory per feature, e.g. image-upload/.
      components/               # UI components specific to this feature.
      api/                      # Communication with the Go backend for this feature.
      types/                    # Shared types for this feature (e.g. API response shapes).
```

- `frontend/` uses a feature-based structure instead of splitting the whole app by technical layer (e.g. a single global `components/` / `lib/api/` / `types/`). As features such as image upload, camera capture, and search results are added, each one's UI, API communication, and types stay colocated under `features/<feature-name>/` instead of being scattered across shared top-level directories.
- `app/` is responsible only for page composition and routing. It may depend on `features/<feature-name>/` for feature-specific UI, API calls, and types, but feature code must not depend on `app/`.
- `features/<feature-name>/` should remain self-contained by default. Feature-specific code should not depend on another feature's internal implementation. Cross-feature or shared dependencies are introduced only when a concrete need arises, rather than being designed in advance.
- Within a feature, `components/`, `api/`, and `types/` are added only as they are actually needed. A feature does not have to contain all three from the start.
- No placeholder directories or `.gitkeep` files are created ahead of need. `features/` and its subdirectories are created only when the first real file for that feature is added.
- State management libraries, data-fetching libraries (e.g. React Query), Server Actions usage, and a shared/common UI component directory are intentionally left undecided here and will be determined in a future Issue when they are actually needed.