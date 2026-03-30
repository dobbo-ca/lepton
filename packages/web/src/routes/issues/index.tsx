import { createResource, Show } from "solid-js";
import { issues } from "../../lib/api";
import IssueRow from "../../components/IssueRow";
import EmptyState from "../../components/EmptyState";
import LoadingSpinner from "../../components/LoadingSpinner";
import TopBar from "../../components/TopBar";

export default function IssuesPage() {
  const [list] = createResource(issues.list);

  return (
    <div class="flex flex-col h-full">
      <TopBar title="Issues" breadcrumbs={[{ label: "Issues" }]} />
      <div class="flex-1 overflow-auto p-6">
        <Show when={list()} fallback={<LoadingSpinner />}>
          {(issueList) => (
            <Show
              when={issueList().length > 0}
              fallback={
                <EmptyState title="No issues" description="Issues from your tracker and native tasks will appear here." />
              }
            >
              <div class="rounded-lg border border-slate-800 overflow-hidden">
                {issueList().map((issue) => (
                  <IssueRow issue={issue} />
                ))}
              </div>
            </Show>
          )}
        </Show>
      </div>
    </div>
  );
}
