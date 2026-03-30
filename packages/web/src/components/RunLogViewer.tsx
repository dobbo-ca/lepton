import { Component } from "solid-js";

interface RunLogViewerProps {
  logs: string;
}

const RunLogViewer: Component<RunLogViewerProps> = (props) => {
  return (
    <div
      class="bg-slate-950 border border-slate-800 rounded-lg overflow-auto max-h-[480px]"
      role="region"
      aria-label="Run logs"
    >
      <pre class="p-4 text-xs text-slate-300 font-mono whitespace-pre-wrap break-all leading-relaxed">
        {props.logs || "No logs available."}
      </pre>
    </div>
  );
};

export default RunLogViewer;
