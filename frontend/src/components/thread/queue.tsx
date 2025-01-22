import { QueueItem } from "@/components/thread/queue-item.tsx";
import { ThreadShape } from "@/db/shapes";
import { Virtuoso } from "react-virtuoso";

export function Queue({
  threads,
  workspaceId,
}: {
  threads: ThreadShape[];
  workspaceId: string;
}) {
  return (
    <div className="divide-y divide-border">
      <Virtuoso
        itemContent={(index) => (
          <QueueItem thread={threads[index]} workspaceId={workspaceId} />
        )}
        totalCount={threads.length}
        useWindowScroll
      />
    </div>
  );
}
