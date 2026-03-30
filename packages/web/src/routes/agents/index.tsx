import { createResource, For, Show } from "solid-js";
import { agents } from "../../lib/api";
import AgentCard from "../../components/AgentCard";
import EmptyState from "../../components/EmptyState";
import LoadingSpinner from "../../components/LoadingSpinner";
import TopBar from "../../components/TopBar";

export default function AgentsPage() {
  const [list, { refetch }] = createResource(agents.list);

  async function handleTrigger(agentId: string) {
    await agents.trigger(agentId);
    refetch();
  }

  return (
    <div class="flex flex-col h-full">
      <TopBar title="Agents" breadcrumbs={[{ label: "Agents" }]} />
      <div class="flex-1 overflow-auto p-6">
        <Show when={list()} fallback={<LoadingSpinner />}>
          {(agentList) => (
            <Show
              when={agentList().length > 0}
              fallback={<EmptyState title="No agents" description="No agents are configured yet." />}
            >
              <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                <For each={agentList()}>
                  {(agent) => <AgentCard agent={agent} onTrigger={handleTrigger} />}
                </For>
              </div>
            </Show>
          )}
        </Show>
      </div>
    </div>
  );
}
