import { ThreadLinkItem } from "@/components/workspace/thread-list-item";
import { Thread } from "@/db/models";
import { Virtuoso } from "react-virtuoso";

export function ThreadList({
  threads,
  workspaceId,
}: {
  threads: Thread[];
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
