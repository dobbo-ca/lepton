import { Component } from "solid-js";

const LoadingSpinner: Component = () => {
  return (
    <div class="flex items-center justify-center py-16" role="status" aria-label="Loading">
      <div
        class="w-6 h-6 border-2 border-slate-700 border-t-slate-400 rounded-full animate-spin"
        aria-hidden="true"
      />
      <span class="sr-only">Loading...</span>
    </div>
  );
};

export default LoadingSpinner;
