import axios from "axios";

const API_BASE_URL =
  (import.meta as ImportMeta & { env?: { VITE_API_BASE_URL?: string } }).env
    ?.VITE_API_BASE_URL ?? "http://localhost:8000";

export type StepStatus =
  | "pending"
  | "running"
  | "completed"
  | "failed"
  | "rollback";

export type WorkflowStep = {
  key: string;
  label: string;
  status: StepStatus;
  message?: string;
  started_at?: string;
  completed_at?: string;
};

export type ProgressState = {
  workflow_id?: string;
  status: "pending" | "running" | "completed" | "failed" | "canceled";
  current_step?: string;
  error?: string;
  steps: WorkflowStep[];
};

export type StartWorkflowPayload = {
  app_name: string;
  namespace: string;
  cluster: string;
};

export async function startWorkflow(payload: StartWorkflowPayload) {
  const res = await axios.post(`${API_BASE_URL}/workflows`, payload);
  return res.data as { workflow_id: string; run_id: string };
}

export async function getWorkflowProgress(workflowId: string) {
  const res = await axios.get(
    `${API_BASE_URL}/workflows/${workflowId}/progress`,
  );
  return res.data as ProgressState;
}

export async function cancelWorkflow(workflowId: string) {
  const res = await axios.post(
    `${API_BASE_URL}/workflows/${workflowId}/cancel`,
  );
  return res.data as { workflow_id: string; status: string };
}
