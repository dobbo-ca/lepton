const API_BASE = import.meta.env.LEPTON_API_URL ?? "http://localhost:8080";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { "Content-Type": "application/json", ...init?.headers },
    ...init,
  });
  if (!res.ok) throw new Error(`API error ${res.status}: ${path}`);
  return res.json() as Promise<T>;
}

// ---- Types ----

export type AgentStatus = "active" | "idle" | "error";

export interface Agent {
  id: string;
  name: string;
  role: string;
  status: AgentStatus;
  lastRunAt: string | null;
  lastRunStatus: RunStatus | null;
}

export type RunStatus = "pending" | "running" | "done" | "failed";

export interface Run {
  id: string;
  agentId: string;
  agentName: string;
  status: RunStatus;
  startedAt: string;
  finishedAt: string | null;
  logs: string;
}

export type IssueStatus = "open" | "in_progress" | "done" | "blocked";
export type IssuePriority = "critical" | "high" | "medium" | "low";

export interface Issue {
  id: string;
  identifier: string;
  title: string;
  status: IssueStatus;
  priority: IssuePriority;
  assigneeId: string | null;
  assigneeName: string | null;
  createdAt: string;
  updatedAt: string;
}

export type RoutineSchedule = string; // cron expression

export interface Routine {
  id: string;
  name: string;
  description: string;
  schedule: RoutineSchedule;
  agentId: string;
  enabled: boolean;
  lastRunAt: string | null;
}

export interface Settings {
  trackerType: string;
  trackerBaseUrl: string;
  k8sConnected: boolean;
}

export interface DashboardData {
  activeAgentCount: number;
  openIssueCount: number;
  recentRuns: Run[];
}

// ---- Agents ----

export const agents = {
  list: () => request<Agent[]>("/api/agents"),
  get: (id: string) => request<Agent>(`/api/agents/${id}`),
  runs: (id: string) => request<Run[]>(`/api/agents/${id}/runs`),
  trigger: (id: string) => request<Run>(`/api/agents/${id}/trigger`, { method: "POST" }),
};

// ---- Issues ----

export const issues = {
  list: () => request<Issue[]>("/api/issues"),
  get: (id: string) => request<Issue>(`/api/issues/${id}`),
};

// ---- Runs ----

export const runs = {
  list: () => request<Run[]>("/api/runs"),
  get: (id: string) => request<Run>(`/api/runs/${id}`),
};

// ---- Routines ----

export const routines = {
  list: () => request<Routine[]>("/api/routines"),
};

// ---- Dashboard ----

export const dashboard = {
  get: () => request<DashboardData>("/api/dashboard"),
};

// ---- Settings ----

export const settings = {
  get: () => request<Settings>("/api/settings"),
  update: (data: Partial<Settings>) =>
    request<Settings>("/api/settings", { method: "PATCH", body: JSON.stringify(data) }),
};
