import { createResource, For, Show } from "solid-js";
import { agents, dashboard } from "../lib/api";
import AgentCard from "../components/AgentCard";
import RunStatusBadge from "../components/RunStatusBadge";
import EmptyState from "../components/EmptyState";
import LoadingSpinner from "../components/LoadingSpinner";
import TopBar from "../components/TopBar";

export default function DashboardPage() {
  const [data, { refetch }] = createResource(dashboard.get);
  const [agentList] = createResource(agents.list);

  async function handleTrigger(agentId: string) {
    await agents.trigger(agentId);
    refetch();
  }

  return (
    <div class="flex flex-col h-full">
      <TopBar title="Dashboard" />
      <div class="flex-1 overflow-auto p-6 space-y-8">
        <Show when={data()} fallback={<LoadingSpinner />}>
          {(d) => (
            <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <SummaryCard label="Active Agents" value={d().activeAgentCount} />
              <SummaryCard label="Open Issues" value={d().openIssueCount} />
              <SummaryCard label="Recent Runs" value={d().recentRuns.length} />
            </div>
          )}
        </Show>

        <section aria-label="Recent runs">
          <h2 class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3">Recent Runs</h2>
          <Show when={data()} fallback={<LoadingSpinner />}>
            {(d) => (
              <Show
                when={d().recentRuns.length > 0}
                fallback={
                  <EmptyState
                    title="No recent runs"
                    description="Runs will appear here once agents are triggered."
                  />
                }
              >
                <div class="rounded-lg border border-slate-800 overflow-hidden">
                  <table class="w-full text-sm">
                    <thead>
                      <tr class="bg-slate-900 text-left">
                        <th class="px-4 py-2.5 text-xs font-medium text-slate-500">Agent</th>
                        <th class="px-4 py-2.5 text-xs font-medium text-slate-500">Status</th>
                        <th class="px-4 py-2.5 text-xs font-medium text-slate-500 hidden sm:table-cell">Started</th>
                        <th class="px-4 py-2.5 text-xs font-medium text-slate-500 hidden sm:table-cell">Finished</th>
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-800">
                      <For each={d().recentRuns}>
                        {(run) => (
                          <tr class="hover:bg-slate-800/50 transition-colors">
                            <td class="px-4 py-3 text-slate-200">{run.agentName}</td>
                            <td class="px-4 py-3">
                              <RunStatusBadge status={run.status} />
                            </td>
                            <td class="px-4 py-3 text-xs text-slate-500 hidden sm:table-cell">
                              {new Date(run.startedAt).toLocaleString()}
                            </td>
                            <td class="px-4 py-3 text-xs text-slate-500 hidden sm:table-cell">
                              {run.finishedAt ? new Date(run.finishedAt).toLocaleString() : "—"}
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

        <section aria-label="Agents">
          <h2 class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3">Agents</h2>
          <Show when={agentList()} fallback={<LoadingSpinner />}>
            {(list) => (
              <Show
                when={list().length > 0}
                fallback={<EmptyState title="No agents" description="No agents are configured yet." />}
              >
                <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                  <For each={list()}>
                    {(agent) => <AgentCard agent={agent} onTrigger={handleTrigger} />}
                  </For>
                </div>
              </Show>
            )}
          </Show>
        </section>
      </div>
    </div>
  );
}

function SummaryCard(props: { label: string; value: number }) {
  return (
    <div class="rounded-lg bg-slate-900 border border-slate-800 p-4 flex flex-col gap-1">
      <span class="text-xs text-slate-500">{props.label}</span>
      <span class="text-2xl font-semibold text-slate-100">{props.value}</span>
    </div>
  );
}
