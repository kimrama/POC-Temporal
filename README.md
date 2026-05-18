# Temporal Progress Mock

Monorepo mock project:

- `frontend`: React + Vite + TypeScript + Tailwind + Axios
- `backend`: Python FastAPI, talks to Temporal using Python SDK
- `workers`: Go Temporal workers
  - `worker-k8s`: hosts the workflow and Kubernetes-like activities
  - `worker-cert`: hosts cert/deploy/verify activities on another task queue
- `docker-compose.yml`: Temporal dev server + UI + backend + frontend + 2 Go workers

## Architecture

```txt
Frontend React
  -> FastAPI backend
    -> Temporal Client
      -> Go worker-k8s
          - ProvisionWorkflow
          - CreateNamespace activity
          - CreateSecret activity
      -> Go worker-cert
          - RequestCertificate activity
          - DeployApplication activity
          - VerifyApplication activity
```

Frontend polls:

```txt
GET /workflows/{workflow_id}/progress
```

Backend queries Temporal workflow:

```txt
query: get_progress
```

## Run

```bash
docker compose up --build
```

Then open:

- Frontend: http://localhost:5173
- Backend API docs: http://localhost:8000/docs
- Temporal UI: http://localhost:8080

## API

Start workflow:

```bash
curl -X POST http://localhost:8000/workflows \
  -H 'Content-Type: application/json' \
  -d '{"app_name":"demo-app","namespace":"demo-ns","cluster":"dev"}'
```

Get progress:

```bash
curl http://localhost:8000/workflows/<workflow_id>/progress
```

Cancel workflow:

```bash
curl -X POST http://localhost:8000/workflows/<workflow_id>/cancel
```

## Notes

This is a mock. Activities only sleep and return fake values.

Important Temporal pattern:

- Workflow keeps progress state in memory.
- State is durable because Temporal records workflow history.
- Backend queries the workflow using `get_progress`.
- Frontend polls backend.
- Side effects are inside activities, not workflow code.
# POC-Temporal
