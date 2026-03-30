import { A } from "@solidjs/router";
import { Component } from "solid-js";
import type { Issue, IssuePriority, IssueStatus } from "../lib/api";

interface IssueRowProps {
  issue: Issue;
}

const statusConfig: Record<IssueStatus, { label: string; classes: string }> = {
  open: { label: "Open", classes: "bg-slate-700 text-slate-300" },
  in_progress: { label: "In Progress", classes: "bg-blue-900 text-blue-300" },
  done: { label: "Done", classes: "bg-emerald-900 text-emerald-300" },
  blocked: { label: "Blocked", classes: "bg-amber-900 text-amber-300" },
};

const priorityConfig: Record<IssuePriority, { label: string; classes: string }> = {
  critical: { label: "Critical", classes: "text-red-400" },
  high: { label: "High", classes: "text-orange-400" },
  medium: { label: "Medium", classes: "text-yellow-400" },
  low: { label: "Low", classes: "text-slate-500" },
};

const IssueRow: Component<IssueRowProps> = (props) => {
  const status = () => statusConfig[props.issue.status] ?? statusConfig.open;
  const priority = () => priorityConfig[props.issue.priority] ?? priorityConfig.medium;

  return (
    <A
      href={`/issues/${props.issue.id}`}
      class="flex items-center gap-4 px-4 py-3 hover:bg-slate-800/60 transition-colors border-b border-slate-800 last:border-0"
      aria-label={`Issue ${props.issue.identifier}: ${props.issue.title}`}
    >
      <span class={`text-xs font-medium shrink-0 ${priority().classes}`} aria-label={`Priority: ${priority().label}`}>
        {priority().label[0]}
      </span>
      <span class="text-xs text-slate-500 font-mono shrink-0 w-16">{props.issue.identifier}</span>
      <span class="text-sm text-slate-200 flex-1 truncate">{props.issue.title}</span>
      {props.issue.assigneeName && (
        <span class="text-xs text-slate-500 shrink-0 hidden sm:block">{props.issue.assigneeName}</span>
      )}
      <span
        class={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium shrink-0 ${status().classes}`}
      >
        {status().label}
      </span>
    </A>
  );
};

export default IssueRow;
