import { Component } from "solid-js";

interface EmptyStateProps {
  title: string;
  description?: string;
}

const EmptyState: Component<EmptyStateProps> = (props) => {
  return (
    <div
      class="flex flex-col items-center justify-center gap-2 py-16 text-center"
      role="status"
      aria-label={props.title}
    >
      <span class="text-3xl text-slate-700" aria-hidden="true">—</span>
      <p class="text-sm font-medium text-slate-400">{props.title}</p>
      {props.description && <p class="text-xs text-slate-600 max-w-xs">{props.description}</p>}
    </div>
  );
};

export default EmptyState;
