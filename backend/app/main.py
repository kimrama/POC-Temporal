import os
import uuid
from typing import Any

from fastapi import FastAPI, HTTPException
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware
from temporalio.client import Client, WorkflowFailureError
from pydantic import BaseModel


TEMPORAL_ADDRESS = os.getenv("TEMPORAL_ADDRESS", "localhost:7233")
TEMPORAL_NAMESPACE = os.getenv("TEMPORAL_NAMESPACE", "default")
WORKFLOW_TASK_QUEUE = os.getenv("WORKFLOW_TASK_QUEUE", "provision-task-queue")
FRONTEND_ORIGIN = os.getenv("FRONTEND_ORIGIN", "http://localhost:5173")


class StartWorkflowRequest(BaseModel):
    app_name: str = "demo-app"
    namespace: str = "demo-ns"
    cluster: str = "dev"


class StartWorkflowResponse(BaseModel):
    workflow_id: str
    run_id: str


app = FastAPI(title="Temporal Progress Mock API")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
    allow_credentials=True,

)


_temporal_client: Client | None = None


async def get_temporal_client() -> Client:
    global _temporal_client
    if _temporal_client is None:
        _temporal_client = await Client.connect(
            TEMPORAL_ADDRESS,
            namespace=TEMPORAL_NAMESPACE,
        )
    return _temporal_client


@app.get("/health")
async def health() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/workflows", response_model=StartWorkflowResponse)
async def start_workflow(payload: StartWorkflowRequest) -> StartWorkflowResponse:
    client = await get_temporal_client()

    workflow_id = f"provision-{payload.app_name}-{uuid.uuid4().hex[:8]}"

    try:
        handle = await client.start_workflow(
            "ProvisionWorkflow",
            {
                "app_name": payload.app_name,
                "namespace": payload.namespace,
                "cluster": payload.cluster,
            },
            id=workflow_id,
            task_queue=WORKFLOW_TASK_QUEUE,
        )
    except Exception as exc:
        raise HTTPException(status_code=500, detail=f"Failed to start workflow: {exc}") from exc

    run_id = handle.result_run_id or handle.first_execution_run_id or ""
    return StartWorkflowResponse(workflow_id=handle.id, run_id=run_id)


@app.get("/workflows/{workflow_id}/progress")
async def get_progress(workflow_id: str) -> dict[str, Any]:
    client = await get_temporal_client()
    handle = client.get_workflow_handle(workflow_id)

    try:
        progress = await handle.query("get_progress")
        print(f"Queried progress for workflow {workflow_id}: {progress}")
        return progress
    except WorkflowFailureError as exc:
        raise HTTPException(status_code=409, detail=str(exc)) from exc
    except Exception as exc:
        msg = str(exc).lower()
        if "consistent query buffer is full" in msg or "resourceexhausted" in msg or "resource exhausted" in msg:
            # Temporal returns ResourceExhausted when the workflow's query buffer is full.
            # Return 429 so callers (frontend) can back off and retry later and include a Retry-After header.
            return JSONResponse(
                status_code=429,
                content={"detail": "Workflow busy, try again later"},
                headers={"Retry-After": "2"},
            )
        raise HTTPException(status_code=404, detail=f"Cannot query workflow: {exc}") from exc


@app.get("/workflows/{workflow_id}/result")
async def get_result(workflow_id: str) -> dict[str, Any]:
    client = await get_temporal_client()
    handle = client.get_workflow_handle(workflow_id)

    try:
        result = await handle.result()
        return {"workflow_id": workflow_id, "result": result}
    except Exception as exc:
        raise HTTPException(status_code=409, detail=str(exc)) from exc


@app.post("/workflows/{workflow_id}/cancel")
async def cancel_workflow(workflow_id: str) -> dict[str, str]:
    client = await get_temporal_client()
    handle = client.get_workflow_handle(workflow_id)
    await handle.cancel()
    return {"workflow_id": workflow_id, "status": "cancel_requested"}
