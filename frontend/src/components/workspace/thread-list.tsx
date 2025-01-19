import { ThreadLinkItem } from "@/components/workspace/thread-list-item";
import { ThreadShape } from "@/db/shapes";
import { Virtuoso } from "react-virtuoso";

export function ThreadList({
  threads,
  workspaceId,
}: {
  threads: ThreadShape[];
  workspaceId: string;
}) {
  return (
    <Virtuoso
      itemContent={(index) => (
        <ThreadLinkItem thread={threads[index]} workspaceId={workspaceId} />
      )}
      totalCount={threads.length}
      useWindowScroll
    />
  );
}
