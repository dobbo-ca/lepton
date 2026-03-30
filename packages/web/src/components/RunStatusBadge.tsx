import { Component } from "solid-js";
import type { RunStatus } from "../lib/api";

interface RunStatusBadgeProps {
  status: RunStatus;
}

const statusConfig: Record<RunStatus, { label: string; classes: string }> = {
  pending: { label: "Pending", classes: "bg-slate-700 text-slate-300" },
  running: { label: "Running", classes: "bg-blue-900 text-blue-300" },
  done: { label: "Done", classes: "bg-emerald-900 text-emerald-300" },
  failed: { label: "Failed", classes: "bg-red-900 text-red-300" },
};

const RunStatusBadge: Component<RunStatusBadgeProps> = (props) => {
  const config = () => statusConfig[props.status] ?? statusConfig.pending;
  return (
    <span
      class={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${config().classes}`}
      aria-label={`Status: ${config().label}`}
    >
      {config().label}
    </span>
  );
};

export default RunStatusBadge;
