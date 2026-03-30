import { Component } from "solid-js";
import type { Agent } from "../lib/api";
import RunStatusBadge from "./RunStatusBadge";

interface AgentCardProps {
  agent: Agent;
  onTrigger?: (agentId: string) => void;
}

const statusDot: Record<Agent["status"], string> = {
  active: "bg-emerald-400",
  idle: "bg-slate-500",
  error: "bg-red-400",
};

const AgentCard: Component<AgentCardProps> = (props) => {
  return (
    <div class="flex flex-col gap-3 p-4 rounded-lg bg-slate-900 border border-slate-800">
      <div class="flex items-start justify-between gap-2">
        <div class="flex flex-col gap-0.5 min-w-0">
          <div class="flex items-center gap-2">
            <span
              class={`w-2 h-2 rounded-full shrink-0 ${statusDot[props.agent.status]}`}
              aria-hidden="true"
            />
            <span class="text-sm font-medium text-slate-100 truncate">{props.agent.name}</span>
          </div>
          <span class="text-xs text-slate-500 pl-4">{props.agent.role}</span>
        </div>
        {props.agent.lastRunStatus && <RunStatusBadge status={props.agent.lastRunStatus} />}
      </div>
      {props.agent.lastRunAt && (
        <p class="text-xs text-slate-500">
          Last run: {new Date(props.agent.lastRunAt).toLocaleString()}
        </p>
      )}
      {props.onTrigger && (
        <button
          type="button"
          onClick={() => props.onTrigger!(props.agent.id)}
          class="mt-auto self-start text-xs px-3 py-1.5 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-slate-100 transition-colors border border-slate-700"
          aria-label={`Trigger agent ${props.agent.name}`}
        >
          Trigger
        </button>
      )}
    </div>
  );
};

export default AgentCard;
