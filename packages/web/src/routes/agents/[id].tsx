import { createResource, For, Show } from "solid-js";
import { useParams } from "@solidjs/router";
import { agents } from "../../lib/api";
import RunStatusBadge from "../../components/RunStatusBadge";
import EmptyState from "../../components/EmptyState";
import LoadingSpinner from "../../components/LoadingSpinner";
import TopBar from "../../components/TopBar";

export default function AgentDetailPage() {
  const params = useParams<{ id: string }>();
  const [agent] = createResource(() => params.id, agents.get);
  const [runs] = createResource(() => params.id, agents.runs);

  return (
    <div class="flex flex-col h-full">
      <Show when={agent()} fallback={<TopBar title="Agent" breadcrumbs={[{ label: "Agents", href: "/agents" }, { label: "Loading…" }]} />}>
        {(a) => (
          <TopBar
            title={a().name}
            breadcrumbs={[{ label: "Agents", href: "/agents" }, { label: a().name }]}
          />
        )}
      </Show>
      <div class="flex-1 overflow-auto p-6 space-y-6">
        <Show when={agent()} fallback={<LoadingSpinner />}>
          {(a) => (
            <div class="rounded-lg bg-slate-900 border border-slate-800 p-5 flex flex-col gap-2">
              <div class="flex items-center gap-3">
                <span class="text-sm font-medium text-slate-100">{a().name}</span>
                <span class="text-xs text-slate-500">{a().role}</span>
              </div>
              <div class="flex items-center gap-2">
                <span class="text-xs text-slate-500">Status:</span>
                <span class="text-xs text-slate-300 capitalize">{a().status}</span>
              </div>
              {a().lastRunAt && (
                <div class="flex items-center gap-2">
                  <span class="text-xs text-slate-500">Last run:</span>
                  <span class="text-xs text-slate-300">{new Date(a().lastRunAt!).toLocaleString()}</span>
                  {a().lastRunStatus && <RunStatusBadge status={a().lastRunStatus!} />}
                </div>
              )}
            </div>
          )}
        </Show>

        <section aria-label="Run history">
          <h2 class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3">Run History</h2>
          <Show when={runs()} fallback={<LoadingSpinner />}>
            {(runList) => (
              <Show
                when={runList().length > 0}
                fallback={<EmptyState title="No runs" description="This agent has no run history yet." />}
              >
                <div class="rounded-lg border border-slate-800 overflow-hidden">
                  <table class="w-full text-sm">
                    <thead>
                      <tr class="bg-slate-900 text-left">
                        <th class="px-4 py-2.5 text-xs font-medium text-slate-500">Status</th>
                        <th class="px-4 py-2.5 text-xs font-medium text-slate-500">Started</th>
                        <th class="px-4 py-2.5 text-xs font-medium text-slate-500 hidden sm:table-cell">Finished</th>
                        <th class="px-4 py-2.5 text-xs font-medium text-slate-500"></th>
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-800">
                      <For each={runList()}>
                        {(run) => (
                          <tr class="hover:bg-slate-800/50 transition-colors">
                            <td class="px-4 py-3">
                              <RunStatusBadge status={run.status} />
                            </td>
                            <td class="px-4 py-3 text-xs text-slate-300">
                              {new Date(run.startedAt).toLocaleString()}
                            </td>
                            <td class="px-4 py-3 text-xs text-slate-500 hidden sm:table-cell">
                              {run.finishedAt ? new Date(run.finishedAt).toLocaleString() : "—"}
                            </td>
                            <td class="px-4 py-3 text-right">
                              <a
                                href={`/runs/${run.id}`}
                                class="text-xs text-slate-400 hover:text-slate-200 transition-colors"
                              >
                                View →
                              </a>
                            </td>
                          </tr>
                        )}
                      </For>
                    </tbody>
                  </table>
                </div>
              </Show>
            )}
          </Show>
        </section>
      </div>
    </div>
  );
}
