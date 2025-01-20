import { ThreadQueueItem } from "@/components/workspace/thread/thread-queue-item";
import { ThreadShape } from "@/db/shapes";
import { Virtuoso } from "react-virtuoso";

export function ThreadQueue({
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
          <ThreadQueueItem thread={threads[index]} workspaceId={workspaceId} />
        )}
        totalCount={threads.length}
        useWindowScroll
      />
    </div>
  );
}
