import { useEffect, useMemo, useState } from "react";
import {
  cancelWorkflow,
  getWorkflowProgress,
  ProgressState,
  startWorkflow,
  StepStatus,
} from "./api";

const defaultProgress: ProgressState = {
  status: "pending",
  steps: [
    { key: "create_namespace", label: "Create namespace", status: "pending" },
    { key: "create_secret", label: "Create secret", status: "pending" },
    {
      key: "request_certificate",
      label: "Request certificate",
      status: "pending",
    },
    {
      key: "deploy_application",
      label: "Deploy application",
      status: "pending",
    },
    {
      key: "verify_application",
      label: "Verify application",
      status: "pending",
    },
  ],
};

const statusIcon: Record<StepStatus, string> = {
  pending: "○",
  running: "⏳",
  completed: "✓",
  failed: "✕",
  rollback: "↩",
};

const statusClass: Record<StepStatus, string> = {
  pending: "bg-slate-100 text-slate-500 border-slate-200",
  running: "bg-blue-50 text-blue-700 border-blue-200",
  completed: "bg-emerald-50 text-emerald-700 border-emerald-200",
  failed: "bg-red-50 text-red-700 border-red-200",
  rollback: "bg-amber-50 text-amber-700 border-amber-200",
};

export function App() {
  const [appName, setAppName] = useState("demo-app");
  const [namespace, setNamespace] = useState("demo-ns");
  const [cluster, setCluster] = useState("dev");
  const [workflowId, setWorkflowId] = useState<string>("");
  const [progress, setProgress] = useState<ProgressState>(defaultProgress);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const completedCount = useMemo(
    () => progress.steps.filter((step) => step.status === "completed").length,
    [progress.steps],
  );

  const progressPercent = Math.round(
    (completedCount / progress.steps.length) * 100,
  );

  const isTerminalStatus =
    progress.status === "completed" ||
    progress.status === "failed" ||
    progress.status === "canceled";

  async function handleStart() {
    setLoading(true);
    setError("");
    setProgress(defaultProgress);

    try {
      const res = await startWorkflow({
        app_name: appName,
        namespace,
        cluster,
      });
      setWorkflowId(res.workflow_id);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to start workflow");
    } finally {
      setLoading(false);
    }
  }

  async function handleCancel() {
    if (!workflowId) return;
    await cancelWorkflow(workflowId);
  }

  useEffect(() => {
    if (!workflowId) return;
    if (isTerminalStatus) return;

    let active = true;
    let timer: number | undefined;

    async function poll() {
      try {
        const latest = await getWorkflowProgress(workflowId);
        if (!active) return;
        setProgress(latest);
        if (
          latest.status === "completed" ||
          latest.status === "failed" ||
          latest.status === "canceled"
        ) {
          if (timer !== undefined) {
            window.clearInterval(timer);
          }
        }
      } catch (err) {
        if (!active) return;
        setError(
          err instanceof Error ? err.message : "Failed to query progress",
        );
      }
    }

    timer = window.setInterval(() => {
      void poll();
    }, 1500);
    poll();

    return () => {
      active = false;
      if (timer !== undefined) {
        window.clearInterval(timer);
      }
    };
  }, [workflowId, isTerminalStatus]);

  const isRunning =
    progress.status === "running" || progress.status === "pending";

  return (
    <main className="min-h-screen p-6 text-slate-900">
      <div className="mx-auto max-w-4xl">
        <section className="mb-6 rounded-3xl bg-white p-6 shadow-sm ring-1 ring-slate-200">
          <h1 className="mt-2 text-3xl font-bold">Temporal Mock</h1>

          <div className="mt-6 grid gap-3 md:grid-cols-3">
            <label className="block">
              <span className="text-sm font-medium text-slate-700">
                App name
              </span>
              <input
                className="mt-1 w-full rounded-xl border border-slate-300 px-3 py-2 outline-none focus:border-blue-500"
                value={appName}
                onChange={(event) => setAppName(event.target.value)}
              />
            </label>
            <label className="block">
              <span className="text-sm font-medium text-slate-700">
                Namespace
              </span>
              <input
                className="mt-1 w-full rounded-xl border border-slate-300 px-3 py-2 outline-none focus:border-blue-500"
                value={namespace}
                onChange={(event) => setNamespace(event.target.value)}
              />
            </label>
            <label className="block">
              <span className="text-sm font-medium text-slate-700">
                Cluster
              </span>
              <input
                className="mt-1 w-full rounded-xl border border-slate-300 px-3 py-2 outline-none focus:border-blue-500"
                value={cluster}
                onChange={(event) => setCluster(event.target.value)}
              />
            </label>
          </div>

          <div className="mt-5 flex flex-wrap gap-3">
            <button
              className="rounded-xl bg-slate-900 px-4 py-2 font-medium text-white disabled:cursor-not-allowed disabled:opacity-50"
              onClick={handleStart}
              disabled={loading || !appName || !namespace || !cluster}
            >
              {loading ? "Starting..." : "Start workflow"}
            </button>

            <button
              className="rounded-xl border border-slate-300 px-4 py-2 font-medium text-slate-700 disabled:cursor-not-allowed disabled:opacity-50"
              onClick={handleCancel}
              disabled={!workflowId || !isRunning}
            >
              Cancel
            </button>
          </div>

          {workflowId && (
            <p className="mt-4 rounded-xl bg-slate-50 px-3 py-2 font-mono text-sm text-slate-700">
              workflow_id: {workflowId}
            </p>
          )}

          {error && (
            <p className="mt-4 rounded-xl bg-red-50 px-3 py-2 text-sm text-red-700">
              {error}
            </p>
          )}
        </section>

        <section className="rounded-3xl bg-white p-6 shadow-sm ring-1 ring-slate-200">
          <div className="flex items-center justify-between gap-3">
            <div>
              <h2 className="text-xl font-semibold">Progress</h2>
              <p className="text-sm text-slate-500">
                status: <span className="font-medium">{progress.status}</span>
                {progress.current_step
                  ? ` · current: ${progress.current_step}`
                  : ""}
              </p>
            </div>
            <div className="text-right text-sm font-medium text-slate-600">
              {progressPercent}%
            </div>
          </div>

          <div className="mt-4 h-3 overflow-hidden rounded-full bg-slate-100">
            <div
              className="h-full rounded-full bg-slate-900 transition-all"
              style={{ width: `${progressPercent}%` }}
            />
          </div>

          <div className="mt-6 space-y-3">
            {progress.steps.map((step) => (
              <div
                key={step.key}
                className={`rounded-2xl border p-4 transition ${statusClass[step.status]}`}
              >
                <div className="flex items-center gap-3">
                  <div className="flex h-9 w-9 items-center justify-center rounded-full bg-white/70 text-lg">
                    {statusIcon[step.status]}
                  </div>
                  <div>
                    <p className="font-semibold">{step.label}</p>
                    <p className="text-sm opacity-80">
                      {step.message || step.status}
                    </p>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {progress.error && (
            <pre className="mt-5 overflow-auto rounded-2xl bg-red-50 p-4 text-sm text-red-800">
              {progress.error}
            </pre>
          )}
        </section>
      </div>
    </main>
  );
}
